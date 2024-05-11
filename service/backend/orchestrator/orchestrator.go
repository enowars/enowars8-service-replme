package orchestrator

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)


type Orchestrator struct {
	Context context.Context
	Client *client.Client
	HostIP string
}

type RunContainerOptions struct {
	ImageTag string
	ContainerName string
	Ports nat.PortMap
}

func Instance() Orchestrator {
	ctx := context.Background()

	opts := []client.Opt{
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	}

	ip := "127.0.0.1"

	ips, err := net.LookupIP("dind")
	if err == nil {
		ip = ips[0].String()
		opts = append(opts, client.WithHost(fmt.Sprintf("https://%s:2376", ip)))
	}
	cli, err := client.NewClientWithOpts(opts...)

	if err != nil {
		panic(err)
	}

	defer cli.Close()

	return Orchestrator{
		Context: ctx,
		Client: cli,
		HostIP: ip,
	}
}

func (o *Orchestrator) BuildImage(path string, tag string) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	ExcludePatterns := []string{}

	exclude, err := os.ReadFile(path + ".dockerignore")

	if err == nil {
		ExcludePatterns = strings.Split(string(exclude), "\n")
	}

	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		ExcludePatterns: ExcludePatterns,
	})

	if err != nil {
		log.Fatal(err, " :unable to create tar")
	}

	opts := types.ImageBuildOptions{
		Dockerfile:  "Dockerfile",
		Tags:        []string { tag },
		Remove:      true,
	}

	res, err := o.Client.ImageBuild(o.Context, tar, opts)

	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)

	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

func (o *Orchestrator) CreateContainer(opts RunContainerOptions) (container.CreateResponse, error) {
	return o.Client.ContainerCreate(
		o.Context, &container.Config { Image: opts.ImageTag },
		&container.HostConfig{
			PortBindings: opts.Ports,
		},
		nil,
		nil,
		opts.ContainerName,
	)
}

func (o *Orchestrator) GetContainer(name string) (*types.Container, *types.ContainerJSON, bool) {
	container, err := o.Client.ContainerList(o.Context, container.ListOptions {
		All: true,
	})

	if err != nil {
		log.Fatal(err, " :could not get container list")
	}

	for _, c := range container {
		for _, v := range c.Names {
			if v[1:] == name {
				information, err := o.Client.ContainerInspect(o.Context, c.ID)
				return &c, &information, (err == nil && information.State.Running)
			}
		}
	}

	return nil, nil, false
}

func (o *Orchestrator) StartContainer(id string) error {
	return o.Client.ContainerStart(o.Context, id, container.StartOptions{ })
}

func (o *Orchestrator) GetContainerAddress(id string) (*string, *uint16, error) {
	container, err := o.Client.ContainerList(o.Context, container.ListOptions {
		All: true,
	})

	if err != nil {
		return nil, nil, err
	}

	for _, c := range container {
		if c.ID == id {

			if len(c.Ports) == 0 {
				return nil, nil, types.ErrorResponse{
					Message: "Container has no exposed port",
				}
			}

			if len(c.NetworkSettings.Networks) == 0 {
				return nil, nil, types.ErrorResponse{
					Message: "Container has no network",
				}
			}

			var ipAddress string

			for _, network := range c.NetworkSettings.Networks {
				ipAddress = network.IPAddress
				break
			}

			return &ipAddress, &c.Ports[0].PublicPort, nil
		}
	}

	return nil, nil, types.ErrorResponse{
		Message: "Container not found",
	}
}

func (o *Orchestrator) EnsureContainerStarted(
	name string,
	username string,
	password string,
) (*string, *uint16, error) {
	var id string

	container, _, running := o.GetContainer(name)

	if container != nil {
		id = container.ID
	} else {
		response, err := o.CreateContainer(RunContainerOptions {
			ImageTag: "ptwhy",
			ContainerName: name,
			Ports: nat.PortMap {
				nat.Port("3000/tcp"): []nat.PortBinding {
					{
						HostIP: o.HostIP,
						HostPort: "0",
					},
				},
			},
		})

		if err != nil {
			return nil, nil, types.ErrorResponse {
				Message: "Container creation failed",
			}
		}

		id = response.ID

	}

	if !running {
		err := o.StartContainer(id)
		if err != nil {
			return nil, nil, types.ErrorResponse {
				Message: "Container start failed",
			}
		}
	}

	ip, port, err := o.GetContainerAddress(id)

	if err != nil {
		return nil, nil, err
	}

	return ip, port, nil
}

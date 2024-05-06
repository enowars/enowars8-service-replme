package orchestrator

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)


type Orchestrator struct {
	Context context.Context
	Client *client.Client
}

type RunContainerOptions struct {
	ImageTag string
	ContainerName string
	Ports nat.PortMap
}

func Instance() Orchestrator {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	return Orchestrator{
		Context: ctx,
		Client: cli,
	}
}

func (o *Orchestrator) BuildImage(path string, tag string) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})

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
				return &c, &information, (err != nil && information.State.Running)
			}
		}
	}

	return nil, nil, false
}

func (o *Orchestrator) StartContainer(id string) error {
	return o.Client.ContainerStart(o.Context, id, container.StartOptions{ })
}

func (o *Orchestrator) GetContainerPort(id string) (*uint16, error) {
	container, err := o.Client.ContainerList(o.Context, container.ListOptions {
		All: true,
	})

	if err != nil {
		return nil, err
	}

	for _, c := range container {
		if c.ID == id {

			if len(c.Ports) == 0 {
				return nil, types.ErrorResponse{
					Message: "Container has no exposed port",
				}
			}

			return &c.Ports[0].PublicPort, nil
		}
	}

	return nil, types.ErrorResponse{
		Message: "Container not found",
	}
}

func (o *Orchestrator) EnsureContainerStarted(
	name string,
	username string,
	password string,
) (*uint16, error) {
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
						HostIP: "127.0.0.1",
						HostPort: "0",
					},
				},
			},
		})

		if err != nil {
			return nil, types.ErrorResponse {
				Message: "Container creation failed",
			}
		}

		id = response.ID

	}

	if !running {
		err := o.StartContainer(id)
		if err != nil {
			return nil, types.ErrorResponse {
				Message: "Container start failed",
			}
		}
	}

	port, err := o.GetContainerPort(id)

	if err != nil {
		return nil, err
	}

	return port, nil
}

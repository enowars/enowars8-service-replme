package service

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"replme/types"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	Context context.Context
	Client  *client.Client
	HostIP  string
	ApiKey  string
}

func Docker(apiKey string) DockerService {
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
		log.Fatal(err)
	}

	defer cli.Close()

	return DockerService{
		Context: ctx,
		Client:  cli,
		HostIP:  ip,
		ApiKey:  apiKey,
	}
}

func (docker *DockerService) BuildImage(imageDir string, tag string) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	ExcludePatterns := []string{}

	exclude, err := os.ReadFile(path.Join(imageDir, ".dockerignore"))

	if err == nil {
		ExcludePatterns = strings.Split(string(exclude), "\n")
	}

	tar, err := archive.TarWithOptions(imageDir, &archive.TarOptions{
		ExcludePatterns: ExcludePatterns,
	})

	if err != nil {
		log.Fatal(err, " :unable to create tar")
	}

	opts := dockerTypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{tag},
		Remove:     true,
		// ForceRemove: true,
		// NoCache:     true,
	}

	res, err := docker.Client.ImageBuild(docker.Context, tar, opts)

	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)

	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

func (docker *DockerService) CreateContainer(
	opts types.RunContainerOptions,
) (*container.CreateResponse, error) {

	volumeEtc, err := docker.Client.VolumeCreate(docker.Context, volume.CreateOptions{
		Name: fmt.Sprintf("%s_etc", opts.ContainerName),
	})

	if err != nil {
		return nil, err
	}

	volumeHome, err := docker.Client.VolumeCreate(docker.Context, volume.CreateOptions{
		Name: fmt.Sprintf("%s_home", opts.ContainerName),
	})

	if err != nil {
		return nil, err
	}

	container, err := docker.Client.ContainerCreate(
		docker.Context,
		&container.Config{
			Image: opts.ImageTag,
			Env:   []string{fmt.Sprintf("API_KEY=%s", docker.ApiKey)},
		},
		&container.HostConfig{
			PortBindings: opts.Ports,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volumeEtc.Name,
					Target: "/etc",
				},
				{
					Type:   mount.TypeVolume,
					Source: volumeHome.Name,
					Target: "/home",
				},
			},
		},
		nil,
		nil,
		opts.ContainerName,
	)

	return &container, err
}

func (docker *DockerService) GetContainers(imageReference string) ([]dockerTypes.Container, error) {

	images, err := docker.Client.ImageList(docker.Context, image.ListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: imageReference,
		}),
	})

	if err != nil || len(images) == 0 {
		return nil, err
	}

	return docker.Client.ContainerList(docker.Context, container.ListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "ancestor",
				Value: images[0].ID,
			},
		),
	})
}

func (docker *DockerService) GetContainer(name string) (
	*dockerTypes.Container,
	*dockerTypes.ContainerJSON,
	bool,
) {
	container, err := docker.Client.ContainerList(docker.Context, container.ListOptions{
		All: true,
	})

	if err != nil {
		log.Fatal(err, " :could not get container list")
	}

	for _, c := range container {
		for _, v := range c.Names {
			if v[1:] == name {
				information, err := docker.Client.ContainerInspect(docker.Context, c.ID)
				return &c, &information, (err == nil && information.State.Running)
			}
		}
	}

	return nil, nil, false
}

func (docker *DockerService) VolumesPrune() (dockerTypes.VolumesPruneReport, error) {
	return docker.Client.VolumesPrune(docker.Context, filters.NewArgs(
		filters.KeyValuePair{
			Key:   "all",
			Value: "1",
		},
	))
}

func (docker *DockerService) StartContainerById(id string) error {
	return docker.Client.ContainerStart(docker.Context, id, container.StartOptions{})
}

func (docker *DockerService) StartContainerByName(name string) {
	c, _, running := docker.GetContainer(name)
	if !running {
		docker.Client.ContainerStart(docker.Context, c.ID, container.StartOptions{})
	}
}

func (docker *DockerService) StopContainerById(id string) error {
	return docker.Client.ContainerStop(docker.Context, id, container.StopOptions{})
}

func (docker *DockerService) StopContainerByName(name string) {
	c, _, running := docker.GetContainer(name)
	if running {
		docker.Client.ContainerStop(docker.Context, c.ID, container.StopOptions{})
	}
}

func (docker *DockerService) KillContainerById(id string) error {
	return docker.Client.ContainerKill(docker.Context, id, "SIGKILL")
}

func (docker *DockerService) KillContainerByName(name string) {
	c, _, running := docker.GetContainer(name)
	if running {
		docker.Client.ContainerKill(docker.Context, c.ID, "SIGKILL")
	}
}

func (docker *DockerService) RemoveContainerById(id string) error {
	return docker.Client.ContainerRemove(docker.Context, id, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

func (docker *DockerService) GetContainerAddress(id string) (*string, *uint16, error) {
	container, err := docker.Client.ContainerList(docker.Context, container.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, nil, err
	}

	for _, c := range container {
		if c.ID == id {

			if len(c.Ports) == 0 {
				return nil, nil, dockerTypes.ErrorResponse{
					Message: "Container has no exposed port",
				}
			}

			if len(c.NetworkSettings.Networks) == 0 {
				return nil, nil, dockerTypes.ErrorResponse{
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

	return nil, nil, dockerTypes.ErrorResponse{
		Message: "Container not found",
	}
}

func (docker *DockerService) EnsureContainerStarted(
	name string,
) (*string, *uint16, error) {
	var id string

	container, _, running := docker.GetContainer(name)

	if container != nil {
		id = container.ID
	} else {
		response, err := docker.CreateContainer(types.RunContainerOptions{
			ImageTag:      "ptwhy",
			ContainerName: name,
			Ports: nat.PortMap{
				nat.Port("3000/tcp"): []nat.PortBinding{
					{
						HostIP:   docker.HostIP,
						HostPort: "0",
					},
				},
			},
		})

		if err != nil {
			fmt.Println(err)
			return nil, nil, dockerTypes.ErrorResponse{
				Message: "Container creation failed",
			}
		}

		id = response.ID

	}

	if !running {
		err := docker.StartContainerById(id)
		if err != nil {
			return nil, nil, dockerTypes.ErrorResponse{
				Message: "Container start failed",
			}
		}
	}

	ip, port, err := docker.GetContainerAddress(id)

	if err != nil {
		return nil, nil, err
	}

	return ip, port, nil
}

func (docker *DockerService) GetContainerPort(name string) *uint16 {

	container, _, running := docker.GetContainer(name)
	if !running {
		return nil
	}
	_, port, _ := docker.GetContainerAddress(container.ID)

	return port
}

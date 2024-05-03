package orchestrator

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/opencontainers/image-spec/specs-go/v1"
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
		&network.NetworkingConfig{},
		&v1.Platform{},
		opts.ContainerName,
	)
}

func (o *Orchestrator) GetContainer(name string) *types.Container {
	container, err := o.Client.ContainerList(o.Context, container.ListOptions {
		All: true,
	})

	if err != nil {
		log.Fatal(err, " :could not get container list")
	}

	for _, c := range container {
		if slices.Contains(c.Names, name) {
			return &c
		}
	}

	return nil
}

func (o *Orchestrator) StartContainer() {
}

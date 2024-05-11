package types

import (
	"github.com/docker/go-connections/nat"
)

type RunContainerOptions struct {
	ImageTag      string
	ContainerName string
	Ports         nat.PortMap
}

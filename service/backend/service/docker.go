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
	"path/filepath"
	"strings"
	"time"

	"replme/types"
	"replme/util"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type DockerService struct {
	Context           context.Context
	Client            *client.Client
	HostIP            string
	ImgPath           string
	ImgTag            string
	ApiKey            string
	ContainerLogsPath string
	MutexMap          util.MutexMap
}

func Docker(apiKey string, imgPath string, imgTag string, containerLogsPath string) DockerService {
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

	err = util.MakeDirIfNotExists(containerLogsPath)

	if err != nil {
		log.Fatal(err)
	}

	return DockerService{
		Context:           ctx,
		Client:            cli,
		HostIP:            ip,
		ImgPath:           imgPath,
		ImgTag:            imgTag,
		ApiKey:            apiKey,
		ContainerLogsPath: containerLogsPath,
		MutexMap:          *util.MutexMapNew(),
	}
}

func (docker *DockerService) BuildImage() {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	ExcludePatterns := []string{}

	exclude, err := os.ReadFile(path.Join(docker.ImgPath, ".dockerignore"))

	if err == nil {
		ExcludePatterns = strings.Split(string(exclude), "\n")
	}

	tar, err := archive.TarWithOptions(docker.ImgPath, &archive.TarOptions{
		ExcludePatterns: ExcludePatterns,
	})

	if err != nil {
		log.Fatal(err, " :unable to create tar")
	}

	httpProxy := "http://proxy.prod.bambi.ovh:3128"
	buildArgs := make(map[string]*string)
	buildArgs["http_proxy"] = &httpProxy
	buildArgs["https_proxy"] = &httpProxy
	buildArgs["HTTP_PROXY"] = &httpProxy
	buildArgs["HTTPS_PROXY"] = &httpProxy
	opts := dockerTypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{docker.ImgTag},
		Remove:     true,
		// ForceRemove: true,
		// NoCache:     true,
		BuildArgs: buildArgs,
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

func (docker *DockerService) CreateReplContainer(
	opts types.RunContainerOptions,
) (*container.CreateResponse, error) {
	util.SLogger.Debugf("[%-25s] Creating container", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]))

	util.SLogger.Debugf("[%-25s] Creating volume \"etc\"", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]))
	start := time.Now()
	volumeEtc, err := docker.Client.VolumeCreate(docker.Context, volume.CreateOptions{
		Name: fmt.Sprintf("%s_etc", opts.ContainerName),
	})
	util.SLogger.Debugf("[%-25s] Creating volume \"etc\" took %v", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]), time.Since(start))

	if err != nil {
		return nil, err
	}

	util.SLogger.Debugf("[%-25s] Creating volume \"home\"", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]))
	start = time.Now()
	volumeHome, err := docker.Client.VolumeCreate(docker.Context, volume.CreateOptions{
		Name: fmt.Sprintf("%s_home", opts.ContainerName),
	})
	util.SLogger.Debugf("[%-25s] Creating volume \"home\" took %v", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]), time.Since(start))

	if err != nil {
		return nil, err
	}

	pidsLimit := int64(256)

	util.SLogger.Debugf("[%-25s] Creating volume \"home\"", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]))
	container, err := docker.Client.ContainerCreate(
		docker.Context,
		&container.Config{
			Image: opts.ImageTag,
			Env: []string{
				fmt.Sprintf("API_KEY=%s", docker.ApiKey),
				"GIN_MODE=release",
			},
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
			Resources: container.Resources{
				CPUQuota:  20000,
				PidsLimit: &pidsLimit,
			},
			// LogConfig: container.LogConfig{
			// 	Type: "none",
			// },
		},
		nil,
		nil,
		opts.ContainerName,
	)

	return &container, err
}

func (docker *DockerService) CreateDevenvContainer(
	devenvPath string,
	mountPath string,
	opts types.RunContainerOptions,
) (*container.CreateResponse, error) {
	util.SLogger.Debugf("[%-25s] Creating container", fmt.Sprintf("NM:%s..", opts.ContainerName[:5]))

	pidsLimit := int64(256)

	container, err := docker.Client.ContainerCreate(
		docker.Context,
		&container.Config{
			Image: opts.ImageTag,
			Env: []string{
				fmt.Sprintf("API_KEY=%s", docker.ApiKey),
				"GIN_MODE=release",
			},
		},
		&container.HostConfig{
			PortBindings: opts.Ports,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: devenvPath,
					Target: mountPath,
				},
			},
			Resources: container.Resources{
				CPUQuota:  20000,
				PidsLimit: &pidsLimit,
			},
			// LogConfig: container.LogConfig{
			// 	Type: "none",
			// },
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
		util.SLogger.Errorf("[%-25s] Failed to get container list: %s", fmt.Sprintf("NM:%s..", name), err.Error())
		return nil, nil, false
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

	out, err := docker.Client.ContainerLogs(docker.Context, c.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: false,
	})

	if err == nil {
		t := time.Now()

		timestamp := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		logFilePath := filepath.Join(docker.ContainerLogsPath, fmt.Sprintf("%s_%s", timestamp, c.ID))
		logFile, err := os.Create(logFilePath)
		if err == nil {
			io.Copy(logFile, out)
			logFile.Close()
		}
		out.Close()
	}

	if running {
		docker.Client.ContainerKill(docker.Context, c.ID, "SIGKILL")
	}
}

func (docker *DockerService) RemoveContainerById(id string) error {

	out, err := docker.Client.ContainerLogs(docker.Context, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: false,
	})

	if err == nil {
		logFilePath := filepath.Join(docker.ContainerLogsPath, id)
		logFile, err := os.Create(logFilePath)
		if err == nil {
			io.Copy(logFile, out)
			logFile.Close()
		}
		out.Close()
	}

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

func (docker *DockerService) EnsureReplContainerStarted(
	name string,
) (*string, *uint16, error) {
	var id string

	container, _, running := docker.GetContainer(name)

	if container != nil {
		id = container.ID
	} else {
		response, err := docker.CreateReplContainer(types.RunContainerOptions{
			ImageTag:      docker.ImgTag,
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

func (docker *DockerService) EnsureDevenvContainerStarted(
	devenvPath string,
	mountPath string,
) (*string, *string, *uint16, error) {
	response, err := docker.CreateDevenvContainer(
		devenvPath,
		mountPath,
		types.RunContainerOptions{
			ImageTag:      docker.ImgTag,
			ContainerName: uuid.NewString(),
			Ports: nat.PortMap{
				nat.Port("3000/tcp"): []nat.PortBinding{
					{
						HostIP:   docker.HostIP,
						HostPort: "0",
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, dockerTypes.ErrorResponse{
			Message: "Container creation failed",
		}
	}
	id := response.ID
	err = docker.StartContainerById(id)
	if err != nil {
		return nil, nil, nil, dockerTypes.ErrorResponse{
			Message: "Container start failed",
		}
	}
	ip, port, err := docker.GetContainerAddress(id)
	if err != nil {
		return nil, nil, nil, err
	}
	return &id, ip, port, nil
}

func (docker *DockerService) GetContainerPort(name string) *uint16 {

	container, _, running := docker.GetContainer(name)
	if !running {
		return nil
	}
	_, port, _ := docker.GetContainerAddress(container.ID)

	return port
}

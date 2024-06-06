package service

import (
	"fmt"
	"time"
)

type CleanupService struct {
	Docker    *DockerService
	ReplState *ReplStateService
}

func Cleanup(docker *DockerService, replState *ReplStateService) CleanupService {
	return CleanupService{
		Docker:    docker,
		ReplState: replState,
	}
}

func (cleanup *CleanupService) DoCleanup() {
	containers, err := cleanup.Docker.GetContainers("ptwhy")

	if err != nil {
		return
	}

	cutoffTime := time.Now().Add(-15 * time.Minute)

	for _, container := range containers {
		created := time.Unix(container.Created, 0)
		if created.Before(cutoffTime) {
			fmt.Println("Removing container: ", container.Names[0])
			cleanup.Docker.RemoveContainerById(container.ID)
			name := container.Names[0][1:] // [1:] because name starts with '/'
			cleanup.ReplState.DeleteContainer(name)
		}
	}

	cleanup.Docker.VolumesPrune()
}

func (cleanup *CleanupService) StartTask() *chan struct{} {

	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Running cleanup ...")
				cleanup.DoCleanup()
				fmt.Println("Cleanup finished")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return &quit
}

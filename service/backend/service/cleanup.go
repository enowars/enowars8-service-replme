package service

import (
	"replme/util"
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
			util.SLogger.Debugf("Removing container: %s", container.Names[0][:10])
			cleanup.Docker.RemoveContainerById(container.ID)
			name := container.Names[0][1:] // [1:] because name starts with '/'
			cleanup.ReplState.DeleteContainer(name)
		}
	}

	util.SLogger.Debug("Pruning volumes starting ..")
	start := time.Now()
	cleanup.Docker.VolumesPrune()
	util.SLogger.Debugf("Pruning volumes [%v]", time.Since(start))
}

func (cleanup *CleanupService) StartTask() *chan struct{} {

	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				util.SLogger.Debug("Cleanup containers starting ..")
				start := time.Now()
				cleanup.DoCleanup()
				util.SLogger.Infof("Cleanup containers [%v]", time.Since(start))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return &quit
}

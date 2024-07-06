package service

import (
	"log"
	"replme/database"
	"replme/model"
	"replme/util"
	"time"
)

type CleanupService struct {
	Docker             *DockerService
	ReplState          *ReplStateService
	ContainerLogsPath  string
	DevenvFilesPath    string
	DevenvFilesTmpPath string
}

func Cleanup(docker *DockerService, replState *ReplStateService, containerLogsPath string, devenvFilesPath string, devenvFilesTmpPath string) CleanupService {

	err := util.MakeDirIfNotExists(containerLogsPath)

	if err != nil {
		log.Fatal(err)
	}

	return CleanupService{
		Docker:             docker,
		ReplState:          replState,
		ContainerLogsPath:  containerLogsPath,
		DevenvFilesPath:    devenvFilesPath,
		DevenvFilesTmpPath: devenvFilesTmpPath,
	}
}

func (cleanup *CleanupService) DoCleanup() {
	containers, err := cleanup.Docker.GetContainers(cleanup.Docker.ImgTag)

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

	util.SLogger.Debug("Cleaning database starting ..")
	start = time.Now()
	database.DB.Unscoped().Where("created_at < ?", cutoffTime).Delete(&model.Devenv{})
	database.DB.Unscoped().Where("created_at < ?", cutoffTime).Delete(&model.User{})
	util.SLogger.Debugf("Cleaning database [%v]", time.Since(start))

	util.SLogger.Debug("Cleaning devenvs starting ..")
	start = time.Now()
	util.DeleteDirsOlderThan(cleanup.DevenvFilesPath, cutoffTime)
	util.DeleteDirsOlderThan(cleanup.DevenvFilesTmpPath, cutoffTime)
	util.SLogger.Debugf("Cleaning devenvs [%v]", time.Since(start))

	util.SLogger.Debug("Cleaning log files starting ...")
	start = time.Now()
	util.DeleteFilesOlderThan(cleanup.ContainerLogsPath, cutoffTime)
	util.SLogger.Debugf("Cleaning log files took [%v]", time.Since(start))

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

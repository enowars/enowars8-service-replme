package service

import (
	"net/http"
	"replme/types"

	"github.com/google/uuid"
)

type ReplStateService struct {
	Map map[string]types.Containers
}

func ReplState() ReplStateService {
	return ReplStateService{
		Map: map[string]types.Containers{},
	}
}

func (repl *ReplStateService) AddContainer(sessionId string, name string, username string, password string) {
	containers := types.Containers{}
	if c, exists := repl.Map[sessionId]; exists {
		containers = c
	} else {
		repl.Map[sessionId] = containers
	}

	if _, exists := containers[name]; !exists {
		containers[name] = types.ContainerData{
			Username: username,
			Password: password,
			Sessions: map[string]http.Cookie{},
		}
	}
}

func (repl *ReplStateService) AddReplSession(sessionId string, name string, cookie http.Cookie) string {

	uid := uuid.New().String()

	containers := types.Containers{}
	if c, exists := repl.Map[sessionId]; exists {
		containers = c
	} else {
		repl.Map[sessionId] = containers
	}

	data := types.ContainerData{
		Sessions: map[string]http.Cookie{},
	}
	if d, exists := containers[name]; exists {
		data = d
	} else {
		containers[name] = data
	}

	data.Sessions[uid] = cookie

	return uid
}

func (repl *ReplStateService) GetContainer(sessionId string, name string) *types.ContainerData {
	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			return &data
		}
	}
	return nil
}

func (repl *ReplStateService) GetContainers(sessionId string) *types.Containers {
	if containers, exists := repl.Map[sessionId]; exists {
		return &containers
	}
	return &types.Containers{}
}

func (repl *ReplStateService) GetReplSession(sessionId string, replUuid string) (*string, *http.Cookie, bool) {

	if containers, exists := repl.Map[sessionId]; exists {
		for name, data := range containers {
			if cookie, exists := data.Sessions[replUuid]; exists {
				return &name, &cookie, true
			}
		}
	}
	return nil, nil, false
}

func (repl *ReplStateService) HasActiveReplSessions(sessionId string, name string) bool {
	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			return len(data.Sessions) > 0
		}
	}
	return false
}

func (repl *ReplStateService) DeleteReplSession(sessionId string, name string, replUuid string) {
	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			delete(data.Sessions, replUuid)
		}
	}
}

func (repl *ReplStateService) DeleteContainer(name string) {
	for _, containers := range repl.Map {
		delete(containers, name)
	}
}

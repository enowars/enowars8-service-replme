package service

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type ContainerCredentials struct {
	Username string
	Password string
}

type ContainerData struct {
	Username string
	Password string
	Sessions map[string] /* ReplUuid -> SessionCookie */ http.Cookie
}

type Containers map[string] /* ContainerName ->*/ ContainerData

type ReplStateService struct {
	Mutex sync.RWMutex
	Map   map[string]Containers
}

func ReplState() ReplStateService {
	return ReplStateService{
		Map: map[string]Containers{},
	}
}

func (repl *ReplStateService) AddContainer(sessionId string, name string, username string, password string) {
	repl.Mutex.Lock()
	defer repl.Mutex.Unlock()

	containers := Containers{}
	if c, exists := repl.Map[sessionId]; exists {
		containers = c
	} else {
		repl.Map[sessionId] = containers
	}

	if _, exists := containers[name]; !exists {
		containers[name] = ContainerData{
			Username: username,
			Password: password,
			Sessions: map[string]http.Cookie{},
		}
	}
}

func (repl *ReplStateService) AddReplSession(sessionId string, name string, cookie http.Cookie) string {
	repl.Mutex.Lock()
	defer repl.Mutex.Unlock()

	uid := uuid.New().String()

	containers := Containers{}
	if c, exists := repl.Map[sessionId]; exists {
		containers = c
	} else {
		repl.Map[sessionId] = containers
	}

	data := ContainerData{
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

func (repl *ReplStateService) GetContainerCredentials(sessionId string, name string) *ContainerCredentials {
	repl.Mutex.RLock()
	defer repl.Mutex.RUnlock()

	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			return &ContainerCredentials{
				Username: data.Username,
				Password: data.Password,
			}
		}
	}

	return nil
}

func (repl *ReplStateService) GetContainerNames(sessionId string) []string {
	repl.Mutex.RLock()
	defer repl.Mutex.RUnlock()

	names := []string{}

	if containers, exists := repl.Map[sessionId]; exists {
		for name := range containers {
			names = append(names, name)
		}
	}

	return names
}

func (repl *ReplStateService) GetReplSession(sessionId string, replUuid string) (*string, *http.Cookie, bool) {
	repl.Mutex.RLock()
	defer repl.Mutex.RUnlock()

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
	repl.Mutex.RLock()
	defer repl.Mutex.RUnlock()

	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			return len(data.Sessions) > 0
		}
	}
	return false
}

func (repl *ReplStateService) DeleteReplSession(sessionId string, name string, replUuid string) {
	repl.Mutex.Lock()
	defer repl.Mutex.Unlock()

	if containers, exists := repl.Map[sessionId]; exists {
		if data, exists := containers[name]; exists {
			delete(data.Sessions, replUuid)
		}
	}
}

func (repl *ReplStateService) DeleteContainer(name string) {
	repl.Mutex.Lock()
	defer repl.Mutex.Unlock()

	for _, containers := range repl.Map {
		delete(containers, name)
	}
}

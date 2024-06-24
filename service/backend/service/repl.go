package service

import (
	"sync"
	"time"
)

type SessionID string
type ContainerName string

type UserSessionData struct {
	Created  time.Time
	Username string
	Password string
}

type UserSessions map[ContainerName]UserSessionData

type ContainerSessionData struct {
	LastAccessed  time.Time
	SessionsCount int
	Mutex         sync.RWMutex
}

type ReplStateService struct {
	UserSessionsMutex      sync.RWMutex
	UserSessionsMap        map[SessionID]UserSessions
	ContainerSessionsMutex sync.RWMutex
	ContainerSessions      map[ContainerName]*ContainerSessionData
}

func ReplState() ReplStateService {
	return ReplStateService{
		UserSessionsMap:   map[SessionID]UserSessions{},
		ContainerSessions: map[ContainerName]*ContainerSessionData{},
	}
}

func (repl *ReplStateService) AddUserSession(sessionId string, name string, username string, password string) {
	sid := SessionID(sessionId)
	cname := ContainerName(name)

	repl.UserSessionsMutex.Lock()
	defer repl.UserSessionsMutex.Unlock()

	sessions := UserSessions{}
	if c, exists := repl.UserSessionsMap[sid]; exists {
		sessions = c
	} else {
		repl.UserSessionsMap[sid] = sessions
	}

	if _, exists := sessions[cname]; !exists {
		sessions[cname] = UserSessionData{
			Created:  time.Now(),
			Username: username,
			Password: password,
		}
	}
}

func (repl *ReplStateService) AddContainerSession(name string) {
	cname := ContainerName(name)

	repl.ContainerSessionsMutex.Lock()
	session, exists := repl.ContainerSessions[cname]
	if !exists {
		session = &ContainerSessionData{}
		repl.ContainerSessions[cname] = session
	}
	repl.ContainerSessionsMutex.Unlock()

	session.Mutex.Lock()
	session.SessionsCount++
	session.LastAccessed = time.Now()
	session.Mutex.Unlock()
}

func (repl *ReplStateService) GetUserSessionData(sessionId string, name string) *UserSessionData {
	repl.UserSessionsMutex.RLock()
	defer repl.UserSessionsMutex.RUnlock()

	sid := SessionID(sessionId)
	cname := ContainerName(name)

	if userSessions, exists := repl.UserSessionsMap[sid]; exists {
		if data, exists := userSessions[cname]; exists {
			return &data
		}
	}

	return nil
}

func (repl *ReplStateService) GetContainerNames(sessionId string) []string {
	sid := SessionID(sessionId)
	names := []string{}

	repl.UserSessionsMutex.RLock()
	defer repl.UserSessionsMutex.RUnlock()

	if userSessions, exists := repl.UserSessionsMap[sid]; exists {
		for name := range userSessions {
			names = append(names, string(name))
		}
	}

	return names
}

func (repl *ReplStateService) ContainerHasActiveSessions(name string) bool {
	repl.UserSessionsMutex.RLock()
	defer repl.UserSessionsMutex.RUnlock()

	cname := ContainerName(name)

	if containerData, exists := repl.ContainerSessions[cname]; exists {
		containerData.Mutex.RLock()
		defer containerData.Mutex.RUnlock()

		return containerData.SessionsCount > 0
	}
	return false
}

func (repl *ReplStateService) DeleteContainerSession(name string, callback func(string)) bool {
	cname := ContainerName(name)

	repl.ContainerSessionsMutex.Lock()
	defer repl.ContainerSessionsMutex.Unlock()

	containerData, exists := repl.ContainerSessions[cname]

	kill := false

	if exists {
		containerData.Mutex.Lock()
		containerData.SessionsCount--
		if containerData.SessionsCount < 1 {
			kill = true
			callback(name)
			delete(repl.ContainerSessions, cname)
		}
		defer containerData.Mutex.Unlock()
	}

	return kill
}

func (repl *ReplStateService) DeleteContainer(name string) {
	cname := ContainerName(name)

	repl.ContainerSessionsMutex.Lock()
	delete(repl.ContainerSessions, cname)
	repl.ContainerSessionsMutex.Unlock()

	repl.UserSessionsMutex.Lock()
	for _, session := range repl.UserSessionsMap {
		if _, exists := session[cname]; exists {
			delete(session, cname)
		}
	}
	repl.UserSessionsMutex.Unlock()
}

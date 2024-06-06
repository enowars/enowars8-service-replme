package types

import "net/http"

type ContainerData struct {
	Username string
	Password string
	Sessions map[string] /* ReplUuid -> SessionCookie */ http.Cookie
}

type Containers map[string] /* ContainerName ->*/ ContainerData

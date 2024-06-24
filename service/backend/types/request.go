package types

import (
	"fmt"
	"net/http"
)

type RequestError struct {
	Code        int
	ContentType string
	Data        []byte
}

func (m *RequestError) Error() string {
	return fmt.Sprintf("RequestError [%d]: %s", m.Code, string(m.Data))
}

type RegisterRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type RequestOptions struct {
	Retries int
	Cookies []*http.Cookie
}

type CreateReplRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type CreateDevenvRequest struct {
	Name     string `json:"name" binding:"required"`
	BuildCmd string `json:"buildCmd"`
	RunCmd   string `json:"runCmd"`
}

type PatchDevenvRequest struct {
	Name     string `json:"name"`
	BuildCmd string `json:"buildCmd"`
	RunCmd   string `json:"runCmd"`
}

type CreateFileRequest struct {
	Name string `json:"name"`
}

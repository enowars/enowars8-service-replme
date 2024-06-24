package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateReplResponse struct {
	ReplUuid string `json:"replUuid"`
}

type AddReplUserResponse struct {
	Id string `json:"id"`
}

type CreateDevenvResponse struct {
	DevenvUuid string `json:"devenvUuid"`
}

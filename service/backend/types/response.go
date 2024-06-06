package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateReplResponse struct {
	ReplUuid string `json:"replUuid"`
}

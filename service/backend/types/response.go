package types

type SuccessResponse struct {
	Success string `json:"success"`
	Port    uint16 `json:"port"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

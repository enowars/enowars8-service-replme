package types

import "fmt"

type ResponseError struct {
	Code    int
	Message string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

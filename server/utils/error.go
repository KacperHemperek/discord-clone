package utils

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d (%s)", e.Code, e.Message)
}

package utils

import "fmt"

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d (%s)", e.Code, e.Message)
}

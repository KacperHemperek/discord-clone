package utils

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d (%s)", e.Code, e.Message)
}

func NewInvalidQueryParamError(paramName string, val any, err error) *APIError {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("%s is an invalid query param, value: %v", paramName, val),
		Cause:   err,
	}
}

func NewNotFoundError(resource, filter string, val any) *APIError {
	return &APIError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s with %s: %v was not found", resource, filter, val),
	}
}

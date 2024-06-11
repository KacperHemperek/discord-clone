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

func NewInvalidQueryParamErr(paramName string, val any, err error) *APIError {
	//if errors.Is(err, store.LimitNumberTooSmallErr) {
	//	return &APIError{
	//		Code:    http.StatusBadRequest,
	//		Message: fmt.Sprintf("%s is too low, number limit query param needs to be at least 1, instead got: %d", paramName, val),
	//		Cause:   err,
	//	}
	//}
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("%s is an invalid query param, value: %v", paramName, val),
		Cause:   err,
	}
}

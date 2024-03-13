package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func ReadBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func HandlerFunc(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logRequest(r, time.Now())
		handlerErr := handler(w, r)
		if handlerErr != nil {
			var err *ApiError
			if errors.As(handlerErr, &err) {
				logApiError(err, r)
				WriteJson(w, err.Code, err)
				return
			}
			logError(handlerErr, r)
			WriteJson(w, http.StatusInternalServerError, &ApiError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server ApiError test restart",
			})
		}
	}
}

func logApiError(err *ApiError, r *http.Request) {
	fmt.Printf("ERROR %s [%s]: %s\n", r.URL, r.Method, err.Error())
	if err.Cause != nil {
		fmt.Printf("CAUSE: %s\n", err.Cause.Error())
	}
}

func logError(err error, r *http.Request) {
	fmt.Printf("ERROR %s [%s]: %s\n", r.URL, r.Method, err.Error())
}

func logRequest(r *http.Request, now time.Time) {
	fmt.Printf("%s [%s] %s\n", r.Method, r.URL, time.Since(now))
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type JSON map[string]interface{}

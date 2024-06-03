package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func ReadBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ReadAndValidateBody(r *http.Request, v interface{}, validate *validator.Validate) error {
	err := ReadBody(r, v)
	if err != nil {
		return err
	}
	return validate.Struct(v)
}

func WsHandler(handler APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logRequest(r, time.Now())

		handlerErr := handler(w, r, &Context{})
		if handlerErr != nil {
			var err *APIError
			if errors.As(handlerErr, &err) {
				logApiError(err, r)
				return
			}
			logError(handlerErr, r)
		}
	}
}

func HandlerFunc(handler APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logRequest(r, time.Now())

		handlerErr := handler(w, r, &Context{})
		if handlerErr != nil {
			var err *APIError
			if errors.As(handlerErr, &err) {
				logApiError(err, r)
				_ = WriteJson(w, err.Code, err)
				return
			}
			logError(handlerErr, r)
			_ = WriteJson(w, http.StatusInternalServerError, &APIError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			})
		}
	}
}

func truncateURL(s *url.URL) string {
	const truncateLength = 40
	if len(s.String()) <= truncateLength {
		return s.String()
	}

	return s.String()[:truncateLength] + "..."
}

func logApiError(err *APIError, r *http.Request) {
	fmt.Printf("ERROR %s [%s]: %s\n", truncateURL(r.URL), r.Method, err.Error())
	if err.Cause != nil {
		fmt.Printf("CAUSE: %s\n", err.Cause.Error())
	}
}

func logError(err error, r *http.Request) {
	fmt.Printf("ERROR %s [%s]: %s\n", truncateURL(r.URL), r.Method, err.Error())
}

func logRequest(r *http.Request, now time.Time) {
	fmt.Printf("%s [%s] %s\n", r.Method, truncateURL(r.URL), time.Since(now))
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

func GetIntParam(r *http.Request, param string) (int, error) {
	params := mux.Vars(r)
	value, ok := params[param]
	if !ok {
		return 0, errors.New("missing param")
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("invalid int param")
	}

	return number, nil
}

type APIHandler func(w http.ResponseWriter, r *http.Request, c *Context) error

type Context struct {
	User *JWTUser
	Conn *websocket.Conn
}

type JSON map[string]interface{}

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

//func AuthFunc(handler Handler) Handler {
//	return func(w http.ResponseWriter, r *http.Request, c *Context) error {
//		authCookie, err := getAuthCookie(r)
//
//		if err != nil {
//			return &ApiError{
//				Code:    http.StatusUnauthorized,
//				Message: "Unauthorized",
//			}
//		}
//
//		if err != nil {
//			return &ApiError{
//				Code:    http.StatusUnauthorized,
//				Message: "Unauthorized",
//			}
//		}
//
//		c.User = user
//
//		return handler(w, r, c)
//	}
//}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

func NewUserCookie(email string) *http.Cookie {
	return &http.Cookie{
		Name:     "auth",
		Value:    email,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
	}
}

func getAuthCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return nil, fmt.Errorf("auth cookie not found")
	}
	return cookie, nil
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type JSON map[string]interface{}

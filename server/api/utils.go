package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/store"
	"net/http"
)

func HandlerFunc(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerErr := handler(w, r, &Context{})
		if handlerErr != nil {
			fmt.Printf("ERROR %s [%s]: %s\n", r.URL, r.Method, handlerErr.Error())
			var err *Error
			if errors.As(handlerErr, &err) {
				WriteJson(w, err.Code, err)
				return
			}
			WriteJson(w, http.StatusInternalServerError, &Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error test restart",
			})
		}
	}
}

//func AuthFunc(handler Handler) Handler {
//	return func(w http.ResponseWriter, r *http.Request, c *Context) error {
//		authCookie, err := getAuthCookie(r)
//
//		if err != nil {
//			return &Error{
//				Code:    http.StatusUnauthorized,
//				Message: "Unauthorized",
//			}
//		}
//
//		if err != nil {
//			return &Error{
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

func ReadBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
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

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d (%s)", e.Code, e.Message)
}

type Handler func(w http.ResponseWriter, r *http.Request, c *Context) error

type Context struct {
	User *store.User
}

type JSON map[string]interface{}

package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/store"
	"log"
	"net/http"
)

type Server struct {
	port int
}

func NewApiServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() {
	router := mux.NewRouter()

	//router.HandleFunc("GET /api/health", HandlerFunc(healthCheckHandler))
	router.HandleFunc("/api/auth", HandlerFunc(AuthFunc(authHandler))).Methods("GET")
	router.HandleFunc("/api/users", HandlerFunc(createUserHandler)).Methods("POST")
	router.HandleFunc("/api/users", HandlerFunc(AuthFunc(func(w http.ResponseWriter, r *http.Request, c *Context) error {
		return WriteJson(w, http.StatusOK, c)
	}))).Methods("GET")

	portStr := fmt.Sprintf(":%d", s.port)

	fmt.Printf("Server is running on port %d\n", s.port)
	log.Fatal(http.ListenAndServe(portStr, router))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, c *Context) error {
	return WriteJson(w, http.StatusOK, &JSON{"status": "OK"})
}

func authHandler(w http.ResponseWriter, r *http.Request, c *Context) error {
	return WriteJson(w, http.StatusOK, c.User)
}

func createUserHandler(w http.ResponseWriter, r *http.Request, c *Context) error {
	user := &store.User{}

	err := ReadBody(r, user)

	if err != nil {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}
	}

	sameEmailUser, _ := store.GetUserByEmail(user.Email)

	if sameEmailUser != nil {
		return &Error{
			Code:    http.StatusConflict,
			Message: "User with this email already exists",
		}
	}

	sameUsernameUser, _ := store.GetUserByUsername(user.Username)

	if sameUsernameUser != nil {
		return &Error{
			Code:    http.StatusConflict,
			Message: "User with this username already exists",
		}
	}

	createdUser, err := store.CreateUser(user.Username, user.Email, user.Password)

	if err != nil {
		return &Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Cause:   err,
		}
	}

	userCookie := NewUserCookie(createdUser.Email)

	http.SetCookie(w, userCookie)

	return WriteJson(w, http.StatusOK, createdUser)
}

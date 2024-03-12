package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/handlers"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
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
	db := store.NewDB()
	store.RunMigrations(db)

	v := validator.New()
	userService := store.NewUserService(db)

	router.HandleFunc("/healthcheck", utils.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"status": "ok"})
	})).Methods("GET")

	router.HandleFunc("/auth/register", utils.HandlerFunc(handlers.NewRegisterUserHandler(&handlers.RegisterUserParams{
		UserService: userService,
		Validator:   v,
	}).Handle)).Methods(http.MethodPost)

	portStr := fmt.Sprintf(":%d", s.port)

	fmt.Printf("Server is running on port %d\n", s.port)
	log.Fatal(http.ListenAndServe(portStr, router))
}

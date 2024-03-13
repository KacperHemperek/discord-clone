package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/handlers"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/rs/cors"
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

	// register all services
	userService := store.NewUserService(db)

	// register all middlewares
	authMiddleware := middlewares.NewAuthMiddleware(&middlewares.AuthMiddlewareParams{UserService: userService})

	router.HandleFunc("/healthcheck", utils.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"status": "ok"})
	})).Methods("GET")

	registerHandler := handlers.NewRegisterUserHandler(&handlers.RegisterUserParams{
		UserService: userService,
		Validator:   v,
	})

	router.HandleFunc(
		"/auth/register",
		utils.HandlerFunc(registerHandler.Handle),
	).Methods(http.MethodPost)

	loginHandler := handlers.NewLoginHandler(&handlers.LoginUserParams{
		UserService: userService,
		Validator:   v,
	})

	router.HandleFunc(
		"/auth/login",
		utils.HandlerFunc(loginHandler.Handle),
	).Methods(http.MethodPost)

	getLoggedInUserHandler := handlers.NewGetLoggedInUserHandler(
		&handlers.GetLoggedInUserParams{
			UserService: userService,
		},
	)
	router.HandleFunc(
		"/auth/me",
		utils.HandlerFunc(authMiddleware.Use(getLoggedInUserHandler.Handle)),
	).Methods(http.MethodGet)

	logoutHandler := handlers.NewLogoutUserHandler()

	router.HandleFunc(
		"/auth/logout",
		utils.HandlerFunc(logoutHandler.Handle),
	).Methods(http.MethodPost)

	portStr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server is running on port %d\n", s.port)

	corsRouter := setupCors(router)

	log.Fatal(http.ListenAndServe(portStr, corsRouter))
}

func setupCors(r *mux.Router) http.Handler {
	acceptedOrigins := []string{"http://localhost:5173", "http://localhost:4201"}
	return cors.New(cors.Options{
		AllowedOrigins:   acceptedOrigins,
		AllowCredentials: true,
	}).Handler(r)
}

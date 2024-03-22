package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/ws"
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

	// register all store services
	userService := store.NewUserService(db)
	friendshipService := store.NewFriendshipService(db)
	chatService := store.NewChatService(db)

	// register all ws services
	notificationsWsService := ws.NewNotificationService()

	// register all middlewares
	authMiddleware := middlewares.NewAuthMiddleware()

	setupRoutes(router, authMiddleware, userService, friendshipService, chatService, notificationsWsService, v)

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

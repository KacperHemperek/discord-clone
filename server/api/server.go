package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/ws"
	"github.com/rs/cors"
	"net/http"
)

type Server struct {
	port int
}

func NewApiServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	db := store.NewDB()
	store.RunMigrations(db)

	v := validator.New()

	// register all store services
	notificationStore := store.NewNotificationService(db, v)
	userService := store.NewUserService(db)
	friendshipService := store.NewFriendshipService(db)
	chatService := store.NewChatService(db)
	messageService := store.NewMessageService(db)

	// register all ws services
	notificationsWsService := ws.NewNotificationService()
	chatWsService := ws.NewChatService()

	// register all middlewares
	authMiddleware := middlewares.NewAuthMiddleware()
	connectWsMiddleware := middlewares.NewConnectWsMiddleware()
	wsAuthMiddleware := middlewares.NewWsAuthMiddleware()
	isChatMemberMiddleware := middlewares.NewIsChatMemberMiddleware(chatService)

	setupRoutes(
		router,
		authMiddleware,
		isChatMemberMiddleware,
		connectWsMiddleware,
		wsAuthMiddleware,
		userService,
		friendshipService,
		chatService,
		messageService,
		notificationStore,
		notificationsWsService,
		chatWsService,
		v,
	)

	portStr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server is running on port %d\n", s.port)

	corsRouter := setupCors(router)

	return http.ListenAndServe(portStr, corsRouter)
}

func setupCors(r *mux.Router) http.Handler {
	acceptedOrigins := []string{"http://localhost:5173", "http://localhost:4201"}
	return cors.New(cors.Options{
		AllowedOrigins:   acceptedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodDelete, http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodPut},
	}).Handler(r)
}

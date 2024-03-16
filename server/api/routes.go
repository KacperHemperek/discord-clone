package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/kacperhemperek/discord-go/handlers"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"net/http"
)

func setupRoutes(
	mux *mux.Router,
	userService *store.UserService,
	friendshipService *store.FriendshipService,
	v *validator.Validate,
	authMiddleware *middlewares.AuthMiddleware,
	notificationsWsService *ws.NotificationService,
) {
	mux.HandleFunc("/healthcheck", utils.HandlerFunc(handlers.HandleHealthcheck())).Methods("GET")

	mux.HandleFunc("/auth/register", utils.HandlerFunc(handlers.HandleRegisterUser(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/login", utils.HandlerFunc(handlers.HandleLogin(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/me", utils.HandlerFunc(authMiddleware.Use(handlers.HandleGetLoggedInUser()))).Methods(http.MethodGet)
	mux.HandleFunc("/auth/logout", utils.HandlerFunc(handlers.HandleLogoutUser())).Methods(http.MethodPost)

	mux.HandleFunc("/friends", utils.HandlerFunc(authMiddleware.Use(handlers.HandleSendFriendRequest(userService, friendshipService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/friends/requests", utils.HandlerFunc(authMiddleware.Use(handlers.HandleGetFriendRequests(userService, friendshipService)))).Methods(http.MethodGet)
	mux.HandleFunc("/friends/requests/{requestId}/accept", utils.HandlerFunc(authMiddleware.Use(handlers.HandleAcceptFriendRequest(userService, friendshipService, v)))).Methods(http.MethodPost)

	mux.HandleFunc("/notifications", utils.HandlerFunc(authMiddleware.Use(handlers.HandleSubscribeNotifications(notificationsWsService)))).Methods(http.MethodGet)
	mux.HandleFunc("/notifications", utils.HandlerFunc(authMiddleware.Use(handlers.HandleCreateNotification(notificationsWsService, v)))).Methods(http.MethodPost)
}

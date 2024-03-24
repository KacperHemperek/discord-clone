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
	authMiddleware middlewares.AuthMiddleware,
	connectWsMiddleware middlewares.ConnectWsMiddleware,
	wsAuthMiddleware middlewares.WsAuthMiddleware,
	userService *store.UserService,
	friendshipService *store.FriendshipService,
	chatService *store.ChatService,
	notificationsWsService *ws.NotificationService,
	v *validator.Validate,
) {

	mux.HandleFunc("/healthcheck", utils.HandlerFunc(handlers.HandleHealthcheck())).Methods(http.MethodGet)

	mux.HandleFunc("/auth/register", utils.HandlerFunc(handlers.HandleRegisterUser(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/login", utils.HandlerFunc(handlers.HandleLogin(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/me", utils.HandlerFunc(authMiddleware(handlers.HandleGetLoggedInUser()))).Methods(http.MethodGet)
	mux.HandleFunc("/auth/logout", utils.HandlerFunc(handlers.HandleLogoutUser())).Methods(http.MethodPost)

	mux.HandleFunc("/friends", utils.HandlerFunc(authMiddleware(handlers.HandleGetFriends(friendshipService)))).Methods(http.MethodGet)
	mux.HandleFunc("/friends", utils.HandlerFunc(authMiddleware(handlers.HandleSendFriendRequest(userService, friendshipService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/friends/{friendID}", utils.HandlerFunc(authMiddleware(handlers.HandleRemoveFriend(friendshipService)))).Methods(http.MethodDelete)
	mux.HandleFunc("/friends/requests", utils.HandlerFunc(authMiddleware(handlers.HandleGetFriendRequests(friendshipService)))).Methods(http.MethodGet)
	mux.HandleFunc("/friends/requests/{requestId}/accept", utils.HandlerFunc(authMiddleware(handlers.HandleAcceptFriendRequest(friendshipService)))).Methods(http.MethodPost)
	mux.HandleFunc("/friends/requests/{requestId}/reject", utils.HandlerFunc(authMiddleware(handlers.HandleRejectFriendRequest(friendshipService)))).Methods(http.MethodPost)

	mux.HandleFunc("/chats", utils.HandlerFunc(authMiddleware(handlers.HandleGetUsersChats(chatService)))).Methods(http.MethodGet)
	mux.HandleFunc("/chats/private", utils.HandlerFunc(authMiddleware(handlers.HandleCreatePrivateChat(chatService, friendshipService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/chats/group", utils.HandlerFunc(authMiddleware(handlers.HandleCreateGroupChat(chatService, userService, v)))).Methods(http.MethodPost)

	mux.HandleFunc(
		"/notifications",
		utils.HandlerFunc(wsAuthMiddleware(handlers.HandleSubscribeNotifications(notificationsWsService))),
	).Methods(http.MethodGet)
	mux.HandleFunc(
		"/notifications",
		utils.HandlerFunc(authMiddleware(handlers.HandleCreateNotification(notificationsWsService, v))),
	).Methods(http.MethodPost)
}

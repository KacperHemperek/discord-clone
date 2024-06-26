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
	isChatMemberMiddleware middlewares.IsChatMemberMiddleware,
	connectWsMiddleware middlewares.ConnectWsMiddleware,
	wsAuthMiddleware middlewares.WsAuthMiddleware,
	userService *store.UserService,
	friendshipService *store.FriendshipService,
	chatService *store.ChatService,
	messageService *store.MessageService,
	notificationStore store.NotificationServiceInterface,
	notificationsWsService *ws.NotificationService,
	chatWsService ws.ChatServiceInterface,
	v *validator.Validate,
) {

	mux.HandleFunc("/healthcheck", utils.HandlerFunc(handlers.HandleHealthcheck())).Methods(http.MethodGet)

	mux.HandleFunc("/auth/register", utils.HandlerFunc(handlers.HandleRegisterUser(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/login", utils.HandlerFunc(handlers.HandleLogin(userService, v))).Methods(http.MethodPost)
	mux.HandleFunc("/auth/me", utils.HandlerFunc(authMiddleware(handlers.HandleGetLoggedInUser()))).Methods(http.MethodGet)
	mux.HandleFunc("/auth/logout", utils.HandlerFunc(handlers.HandleLogoutUser())).Methods(http.MethodPost)

	mux.HandleFunc("/friends", utils.HandlerFunc(authMiddleware(handlers.HandleGetFriends(friendshipService)))).Methods(http.MethodGet)
	mux.HandleFunc("/friends", utils.HandlerFunc(authMiddleware(handlers.HandleSendFriendRequest(userService, notificationStore, notificationsWsService, friendshipService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/friends/{friendID}", utils.HandlerFunc(authMiddleware(handlers.HandleRemoveFriend(friendshipService)))).Methods(http.MethodDelete)
	mux.HandleFunc("/friends/requests", utils.HandlerFunc(authMiddleware(handlers.HandleGetFriendRequests(friendshipService)))).Methods(http.MethodGet)
	mux.HandleFunc("/friends/requests/{requestId}/accept", utils.HandlerFunc(authMiddleware(handlers.HandleAcceptFriendRequest(friendshipService)))).Methods(http.MethodPost)
	mux.HandleFunc("/friends/requests/{requestId}/reject", utils.HandlerFunc(authMiddleware(handlers.HandleRejectFriendRequest(friendshipService)))).Methods(http.MethodPost)

	mux.HandleFunc("/chats", utils.HandlerFunc(authMiddleware(handlers.HandleGetUsersChats(chatService)))).Methods(http.MethodGet)
	mux.HandleFunc("/chats/private", utils.HandlerFunc(authMiddleware(handlers.HandleCreatePrivateChat(chatService, friendshipService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/chats/group", utils.HandlerFunc(authMiddleware(handlers.HandleCreateGroupChat(chatService, userService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/chats/{chatID}/messages", utils.HandlerFunc(authMiddleware(handlers.HandleSendMessage(chatService, messageService, chatWsService, notificationStore, notificationsWsService, v)))).Methods(http.MethodPost)
	mux.HandleFunc("/chats/{chatID}", utils.HandlerFunc(authMiddleware(handlers.HandleGetChatWithMessages(chatService)))).Methods(http.MethodGet)
	mux.HandleFunc("/chats/{chatID}/update-name", utils.HandlerFunc(authMiddleware(handlers.HandleUpdateChatName(chatService, chatWsService, v)))).Methods(http.MethodPut)
	mux.HandleFunc("/chats/{chatID}/members/add", utils.HandlerFunc(authMiddleware(isChatMemberMiddleware(handlers.HandleAddUsersToChat(chatService, v))))).Methods(http.MethodPost)

	mux.HandleFunc("/ws/chats/{chatID}", utils.WsHandler(wsAuthMiddleware(handlers.HandleConnectToChat(chatWsService)))).Methods(http.MethodGet)

	mux.HandleFunc(
		"/ws/notifications",
		utils.WsHandler(wsAuthMiddleware(handlers.HandleSubscribeNotifications(notificationsWsService))),
	).Methods(http.MethodGet)

	mux.HandleFunc("/notifications/friend-requests/mark-as-seen", utils.HandlerFunc(authMiddleware(handlers.HandleMarkFriendRequestNotificationsAsSeen(notificationStore)))).Methods(http.MethodPut)
	mux.HandleFunc("/notifications/new-messages/mark-as-seen", utils.HandlerFunc(authMiddleware(handlers.HandleMarkNewMessageNotificationsAsSeen(notificationStore, chatService, v)))).Methods(http.MethodPut)
	mux.HandleFunc("/notifications/new-messages", utils.HandlerFunc(authMiddleware(handlers.HandleGetNewMessageNotifications(notificationStore)))).Methods(http.MethodGet)
	mux.HandleFunc("/notifications/friend-requests", utils.HandlerFunc(authMiddleware(handlers.HandleGetFriendRequestNotifications(notificationStore)))).Methods(http.MethodGet)
}

package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"net/http"
)

type subscribeNotificationsHandler struct {
	wsNotificationService ws.NotificationServiceInterface
}

type createNotificationHandler struct {
	wsNotificationService ws.NotificationServiceInterface
	validate              *validator.Validate
}

func HandleSubscribeNotifications(notificationWsService ws.NotificationServiceInterface) middlewares.HandlerWithUser {
	handler := &subscribeNotificationsHandler{
		wsNotificationService: notificationWsService,
	}

	return handler.handle
}

func HandleCreateNotification(notificationWsService *ws.NotificationService, validate *validator.Validate) middlewares.HandlerWithUser {
	handler := &createNotificationHandler{
		wsNotificationService: notificationWsService,
		validate:              validate,
	}

	return handler.handle
}

func (h *subscribeNotificationsHandler) handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	h.wsNotificationService.AddConn(user.ID, conn)

	for {
		mType, msg, err := conn.ReadMessage()

		if err != nil {
			return err
		}

		switch mType {
		case websocket.TextMessage:
			fmt.Println("received message: ", string(msg))
		case websocket.CloseMessage:
			return nil
		}
	}

}

func (h *createNotificationHandler) handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
	body := &CreateNotificationBody{}
	if err := utils.ReadBody(r, body); err != nil {
		return err
	}
	if err := h.validate.Struct(body); err != nil {
		return err
	}

	if err := h.wsNotificationService.Notify(user.ID, body.Message); err != nil {
		fmt.Printf("error when notifying user %d: %s", user.ID, err)
	}
	return utils.WriteJson(w, http.StatusCreated, utils.JSON{"message": "notification sent"})
}

type CreateNotificationBody struct {
	Message string `json:"message" validate:"required"`
}

type NewCreateNotificationParams struct {
	WsNotificationService *ws.NotificationService
	Validate              *validator.Validate
}

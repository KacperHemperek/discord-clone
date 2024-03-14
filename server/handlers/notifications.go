package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"net/http"
)

type SubscribeNotificationsHandler struct {
	wsNotificationService *ws.NotificationService
}

type CreateNotificationHandler struct {
	wsNotificationService *ws.NotificationService
	validate              *validator.Validate
}

func (h *SubscribeNotificationsHandler) Handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
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

func (h *CreateNotificationHandler) Handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
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

type NewSubscribeNotificationsParams struct {
	WsNotificationService *ws.NotificationService
}

type NewCreateNotificationParams struct {
	WsNotificationService *ws.NotificationService
	V                     *validator.Validate
}

func NewSubscribeNotificationsHandler(params *NewSubscribeNotificationsParams) *SubscribeNotificationsHandler {
	return &SubscribeNotificationsHandler{
		wsNotificationService: params.WsNotificationService,
	}
}

func NewCreateNotificationHandler(params *NewCreateNotificationParams) *CreateNotificationHandler {
	return &CreateNotificationHandler{
		wsNotificationService: params.WsNotificationService,
		validate:              params.V,
	}
}

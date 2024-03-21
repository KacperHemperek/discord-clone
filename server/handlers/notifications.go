package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"net/http"
)

func HandleSubscribeNotifications(notificationWsService ws.NotificationServiceInterface) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}

		notificationWsService.AddConn(c.User.ID, conn)

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
}

func HandleCreateNotification(notificationWsService *ws.NotificationService, validate *validator.Validate) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		body := &CreateNotificationBody{}
		if err := utils.ReadBody(r, body); err != nil {
			return err
		}
		if err := validate.Struct(body); err != nil {
			return err
		}

		if err := notificationWsService.Notify(c.User.ID, body.Message); err != nil {
			fmt.Printf("error when notifying user %d: %s", c.User.ID, err)
		}
		return utils.WriteJson(w, http.StatusCreated, utils.JSON{"message": "notification sent"})
	}
}

type CreateNotificationBody struct {
	Message string `json:"message" validate:"required"`
}

type NewCreateNotificationParams struct {
	WsNotificationService *ws.NotificationService
	Validate              *validator.Validate
}

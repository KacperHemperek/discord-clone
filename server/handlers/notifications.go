package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type SubscribeNotificationsHandler struct {
	wsNotificationService *WsNotificationService
}

type CreateNotificationHandler struct {
	wsNotificationService *WsNotificationService
	validate              *validator.Validate
}

func (h *SubscribeNotificationsHandler) Handle(w http.ResponseWriter, r *http.Request, user *models.User) error {
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

func (h *CreateNotificationHandler) Handle(w http.ResponseWriter, r *http.Request, user *models.User) error {
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
	NS *WsNotificationService
}

type NewCreateNotificationParams struct {
	NS *WsNotificationService
	V  *validator.Validate
}

func NewSubscribeNotificationsHandler(params *NewSubscribeNotificationsParams) *SubscribeNotificationsHandler {
	return &SubscribeNotificationsHandler{
		wsNotificationService: params.NS,
	}
}

func NewCreateNotificationHandler(params *NewCreateNotificationParams) *CreateNotificationHandler {
	return &CreateNotificationHandler{
		wsNotificationService: params.NS,
		validate:              params.V,
	}
}

type WsNotificationService struct {
	conns map[int]*websocket.Conn
}

func (s *WsNotificationService) AddConn(userId int, conn *websocket.Conn) {
	s.conns[userId] = conn
}

func (s *WsNotificationService) RemoveConn(userId int) {
	delete(s.conns, userId)
}

func (s *WsNotificationService) Notify(userId int, msg string) error {
	conn, ok := s.conns[userId]
	if !ok {
		return fmt.Errorf("no connection for user %d", userId)
	}

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func NewWsNotificationService() *WsNotificationService {
	return &WsNotificationService{
		conns: make(map[int]*websocket.Conn),
	}
}

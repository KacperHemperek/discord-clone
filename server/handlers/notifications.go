package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"log/slog"
	"net/http"
)

func HandleSubscribeNotifications(notificationWsService ws.NotificationServiceInterface) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		connID := notificationWsService.AddConn(c.User.ID, c.Conn)
		for {
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				break
			}
		}
		slog.Info("closing user notification connection", "userID", c.User.ID)
		return notificationWsService.RemoveConn(c.User.ID, connID)
	}
}
func HandleMakeNotificationsSeen(notificationsStore store.NotificationServiceInterface) utils.APIHandler {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		typeFilter := r.URL.Query().Get("type")

		err := notificationsStore.MarkUsersNotificationsAsSeenByType(c.User.ID, typeFilter)

		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &response{
			Message: "notifications marked as seen",
		})
	}
}

type CreateNotificationBody struct {
	Message string `json:"message" validate:"required"`
}

type NewCreateNotificationParams struct {
	WsNotificationService *ws.NotificationService
	Validate              *validator.Validate
}

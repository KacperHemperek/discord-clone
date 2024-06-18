package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/types"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"log/slog"
	"net/http"
)

func HandleSubscribeNotifications(notificationWsService ws.NotificationServiceInterface) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
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
func HandleMarkFriendRequestNotificationsAsSeen(notificationsStore store.NotificationServiceInterface) utils.APIHandler {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		frn := types.FriendRequestNotification
		err := notificationsStore.MarkUsersNotificationsAsSeen(c.User.ID, frn.String())

		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &response{
			Message: "notifications marked as seen",
		})
	}
}

func HandleMarkNewMessageNotificationsAsSeen(notificationsStore store.NotificationServiceInterface, chatsStore store.ChatServiceInterface, validate *validator.Validate) utils.APIHandler {
	type request struct {
		ChatID int `json:"chatId" validate:"number,min=1"`
	}

	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		body := &request{}
		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "provided payload is invalid",
				Cause:   err,
			}
		}

		_, err := chatsStore.GetChatByID(body.ChatID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundError("chat", "id", body.ChatID)
		}
		err = notificationsStore.MarkUsersNewMessageNotificationsAsSeenByChatID(c.User.ID, body.ChatID)

		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &response{
			Message: "notifications marked as seen",
		})
	}
}

func HandleGetNewMessageNotifications(notificationsStore store.NotificationServiceInterface) utils.APIHandler {
	type response struct {
		Notifications []*models.NewMessageNotification `json:"notifications"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		seen := r.URL.Query().Get("seen")
		limit := r.URL.Query().Get("limit")

		seenFilter, err := store.NewBoolFilter(seen)

		if err != nil {
			return utils.NewInvalidQueryParamError("seen", seen, err)
		}

		limitFilter, err := store.NewLimitFilter(limit)

		if err != nil {
			return utils.NewInvalidQueryParamError("limit", limit, err)
		}

		notifications, err := notificationsStore.GetUserNewMessageNotifications(c.User.ID, seenFilter, limitFilter)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &response{
			Notifications: notifications,
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

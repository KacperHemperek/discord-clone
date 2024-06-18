package middlewares

import (
	"database/sql"
	"errors"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
	"slices"
)

type IsChatMemberMiddleware = func(h utils.APIHandler) utils.APIHandler

func NewIsChatMemberMiddleware(chatsStore store.ChatServiceInterface) IsChatMemberMiddleware {
	return func(h utils.APIHandler) utils.APIHandler {
		return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
			chatID, err := utils.GetIntParam(r, "chatID")
			if err != nil {
				return err
			}

			chat, err := chatsStore.GetChatByID(chatID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return utils.NewNotFoundError("chat", "id", chatID)
				}
				return err
			}
			users, err := chatsStore.GetChatMembers(chat.ID)
			userIsMember := slices.ContainsFunc(users, func(user *models.User) bool {
				return user.ID == c.User.ID
			})
			if !userIsMember {
				return &utils.APIError{
					Code:    http.StatusForbidden,
					Message: "User is not a chat member",
				}
			}
			return h(w, r, c)
		}
	}
}

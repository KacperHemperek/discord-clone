package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
	"slices"
	"strings"
)

type CreateChatRequestBody struct {
	UserID int `json:"userId" validation:"required,gte=1"`
}

func HandleCreatePrivateChat(
	chatService store.ChatServiceInterface,
	friendService store.FriendshipServiceInterface,
	validate *validator.Validate,
) utils.APIHandler {
	type response struct {
		Message string `json:"message"`
		ChatID  int    `json:"chatId"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		body := &CreateChatRequestBody{}
		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Could not read body",
				Cause:   err,
			}
		}
		if body.UserID == c.User.ID {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Cannot create private chat with yourself",
			}
		}
		chat, err := chatService.GetPrivateChatByUserIDs(c.User.ID, body.UserID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if chat != nil {
			res := &response{
				Message: "Chat with that user already exists",
				ChatID:  chat.ID,
			}
			return utils.WriteJson(w, http.StatusOK, res)
		}
		_, err = friendService.GetFriendshipByUsers(body.UserID, c.User.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &utils.APIError{
					Code:    http.StatusForbidden,
					Message: "Users are not friends",
				}
			}
			return err
		}
		createdChat, err := chatService.CreatePrivateChatWithUsers(body.UserID, c.User.ID)
		if err != nil {
			return err
		}
		return utils.WriteJson(w, http.StatusOK, &response{
			Message: "Chat created",
			ChatID:  createdChat.ID,
		})
	}
}

type CreateGroupChatRequestBody struct {
	UserIDs []int `json:"userIds"`
}

func HandleCreateGroupChat(
	chatService store.ChatServiceInterface,
	userService store.UserServiceInterface,
	validate *validator.Validate,
) utils.APIHandler {
	type response struct {
		ChatID  int    `json:"chatId"`
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		body := &CreateGroupChatRequestBody{}
		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Request body is not valid",
				Cause:   err,
			}
		}
		if slices.Contains(body.UserIDs, c.User.ID) {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "userIds cannot contain logged in user id",
				Cause:   nil,
			}
		}
		allIDs := append(body.UserIDs, c.User.ID)
		users, err := userService.GetUsersByIDs(allIDs)
		if err != nil {
			return err
		}
		if len(users) != len(allIDs) {
			return &utils.APIError{
				Message: "Not every user exists from provided list",
				Code:    http.StatusNotFound,
			}
		}
		usernames := make([]string, len(users))
		for i, user := range users {
			usernames[i] = user.Username
		}
		chatName := strings.Join(usernames, ", ")
		chat, err := chatService.CreateGroupChat(chatName, allIDs)
		if err != nil {
			return err
		}
		return utils.WriteJson(w, http.StatusCreated, &response{ChatID: chat.ID, Message: "Chat created successfully"})
	}
}

func HandleGetUsersChats(
	chatService store.ChatServiceInterface,
) utils.APIHandler {
	type response struct {
		Chats []*models.ChatWithMembers `json:"chats"`
	}

	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		chats, err := chatService.GetUsersChatsWithMembers(c.User.ID)
		if err != nil {
			return err
		}
		for _, chat := range chats {
			if chat.Type == "private" {
				var otherMember *models.User
				for _, m := range chat.Members {
					if m.ID != c.User.ID {
						otherMember = m
						break
					}
				}
				chat.Name = otherMember.Username
			}
		}
		return utils.WriteJson(w, http.StatusOK, &response{Chats: chats})
	}
}

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
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
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
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
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

	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		chats, err := chatService.GetUsersChatsWithMembers(c.User.ID)
		if err != nil {
			return err
		}
		for _, chat := range chats {
			if chat.Type.Is(types.PrivateChat) {
				newChatName, err := getPrivChatName(c.User.ID, chat.Members)
				if err != nil {
					return err
				}
				chat.Name = newChatName
			}
		}
		return utils.WriteJson(w, http.StatusOK, &response{Chats: chats})
	}
}

type SendMessageRequestBody struct {
	Text string `json:"text"`
}

func HandleSendMessage(
	chatService store.ChatServiceInterface,
	messageService store.MessageServiceInterface,
	chatWsService ws.ChatServiceInterface,
	notificationStore store.NotificationServiceInterface,
	notificationService ws.NotificationServiceInterface,
	validate *validator.Validate,
) utils.APIHandler {
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		chatID, err := utils.GetIntParam(r, "chatID")
		if err != nil {
			return err
		}
		body := &SendMessageRequestBody{}
		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Request body is not valid",
				Cause:   err,
			}
		}
		chat, err := chatService.GetChatByID(chatID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &utils.APIError{
					Code:    http.StatusNotFound,
					Message: "Chat was not found",
					Cause:   err,
				}
			}
			return err
		}
		m, err := messageService.CreateMessageInChat(chat.ID, c.User.ID, body.Text)
		if err != nil {
			return nil
		}
		mwu, err := messageService.EnrichMessageWithUser(m)
		if err != nil {
			return err
		}
		_, err = chatWsService.GetActiveUserIDs(chatID)
		if err != nil {
			return err
		}
		err = chatWsService.BroadcastNewMessage(chatID, mwu)
		if err != nil {
			return err
		}

		activeMemberIDs, err := chatWsService.GetActiveUserIDs(chatID)

		if err != nil {
			return err
		}

		members, err := chatService.GetChatMembersExcluding(chatID, activeMemberIDs)

		if err != nil {
			return err
		}
		memberIDs := make([]int, 0)
		for _, m := range members {
			memberIDs = append(memberIDs, m.ID)
		}

		notifications, err := notificationStore.CreateNewMessageNotificationsForUsers(
			memberIDs,
			&models.NewMessageNotificationData{
				ChatID: chatID,
			},
		)

		if err != nil {
			return err
		}

		for _, n := range notifications {
			err := notificationService.SendNotification(n.UserID, n)
			if err != nil {
				slog.Error("could not send new message notification", "userID", n.UserID)
			}
		}

		return utils.WriteJson(w, http.StatusCreated, &response{
			Message: "Message created successfully",
		})
	}
}

func HandleGetChatWithMessages(
	chatService store.ChatServiceInterface,
) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		chatID, err := utils.GetIntParam(r, "chatID")
		if err != nil {
			return err
		}
		chat, err := chatService.GetChatByID(chatID)
		if err != nil {
			return err
		}
		cwm, err := chatService.EnrichChatWithMessages(chat)
		if err != nil {
			return err
		}
		if cwm.Type.Is(types.PrivateChat) {
			members, err := chatService.GetChatMembers(chat.ID)
			if err != nil {
				return err
			}
			chatName, err := getPrivChatName(c.User.ID, members)
			if err != nil {
				return err
			}
			cwm.Name = chatName
		}
		return utils.WriteJson(w, http.StatusOK, cwm)
	}
}

func HandleConnectToChat(wsChatService ws.ChatServiceInterface) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		chatID, err := utils.GetIntParam(r, "chatID")
		if err != nil {
			return err
		}
		connID := wsChatService.AddChatConn(chatID, c.User.ID, c.Conn)
		for {
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				break
			}
		}
		return wsChatService.CloseConn(chatID, connID)
	}
}

func HandleUpdateChatName(chatService store.ChatServiceInterface, chatWsService ws.ChatServiceInterface, validate *validator.Validate) utils.APIHandler {
	type request struct {
		NewName string `json:"newName" validate:"min=6,max=32"`
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		chatID, err := utils.GetIntParam(r, "chatID")
		if err != nil {
			return err
		}
		body := &request{}
		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Request body is on valid",
			}
		}
		chat, err := chatService.GetChatByID(chatID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &utils.APIError{
					Code:    http.StatusNotFound,
					Message: "Chat not found",
				}
			}
			return err
		}
		if chat.Type.Is(types.PrivateChat) {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "You cannot change name of private chat",
			}
		}
		users, err := chatService.GetChatMembers(chat.ID)
		userIsMember := slices.ContainsFunc(users, func(user *models.User) bool {
			return user.ID == c.User.ID
		})
		if !userIsMember {
			return &utils.APIError{
				Code:    http.StatusForbidden,
				Message: "User can't update name of chat he is not a part of",
			}
		}
		err = chatService.UpdateChatName(chatID, body.NewName)
		if err != nil {
			return err
		}
		err = chatWsService.BroadcastNewChatName(chatID, body.NewName)
		if err != nil {
			return err
		}
		return utils.WriteJson(w, http.StatusOK, &response{Message: "Chat name updated successfully"})
	}
}

func getPrivChatName(loggedInUserID int, members []*models.User) (string, error) {
	var chatName string
	found := false
	for _, m := range members {
		if m.ID != loggedInUserID {
			chatName = m.Username
			found = true
			break
		}
	}
	if found {
		return chatName, nil
	}
	return "", errOtherUserNotFount
}

var errOtherUserNotFount = errors.New("other user then logged in not found")

package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type sendFriendRequestHandler struct {
	validate          *validator.Validate
	userService       *store.UserService
	friendshipService *store.FriendshipService
}

type getFriendRequestsHandler struct {
	userService       *store.UserService
	friendshipService *store.FriendshipService
}

func HandleSendFriendRequest(userService *store.UserService, friendshipService *store.FriendshipService, validate *validator.Validate) middlewares.HandlerWithUser {
	handler := &sendFriendRequestHandler{
		validate:          validate,
		userService:       userService,
		friendshipService: friendshipService,
	}

	return handler.handle
}

func HandleGetFriendRequests(userService *store.UserService, friendshipService *store.FriendshipService) middlewares.HandlerWithUser {
	handler := &getFriendRequestsHandler{
		userService:       userService,
		friendshipService: friendshipService,
	}

	return handler.handle
}

func (h *sendFriendRequestHandler) handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
	body := &SendFriendRequestBody{}

	if err := utils.ReadBody(r, body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
	}

	if err := h.validate.Struct(body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
	}

	userToSendRequest, err := h.userService.FindUserByEmail(body.Email)

	if err != nil {
		if errors.Is(err, store.UserNotFoundError) {
			return &utils.ApiError{Code: http.StatusNotFound, Message: "User with that email not found", Cause: err}
		}
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when searching for user", Cause: err}
	}

	if userToSendRequest.ID == user.ID {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Cannot send friend request to yourself", Cause: nil}
	}

	if err := h.friendshipService.SendFriendRequest(user.ID, userToSendRequest.ID); err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when sending friend request", Cause: err}
	}

	return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
}

func (h *getFriendRequestsHandler) handle(w http.ResponseWriter, _ *http.Request, user *utils.JWTUser) error {
	friendRequests, err := h.friendshipService.GetFriendRequests(user.ID)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend requests", Cause: err}
	}

	return utils.WriteJson(w, http.StatusOK, &utils.JSON{"requests": friendRequests})
}

type SendFriendRequestBody struct {
	Email string `json:"email" validate:"required"`
}

package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type SendFriendRequestBody struct {
	Email string `json:"email" validate:"required"`
}

type SendFriendRequestHandler struct {
	validate    *validator.Validate
	userService *store.UserService
}

func (h *SendFriendRequestHandler) Handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
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

	if err := h.userService.SendFriendRequest(user.ID, userToSendRequest.ID); err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when sending friend request", Cause: err}
	}

	return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
}

type NewSendFriendRequestProps struct {
	Validate    *validator.Validate
	UserService *store.UserService
}

func NewSendFriendRequestHandler(props *NewSendFriendRequestProps) *SendFriendRequestHandler {
	return &SendFriendRequestHandler{
		validate:    props.Validate,
		userService: props.UserService,
	}
}

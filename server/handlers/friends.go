package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
	"strconv"
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

type acceptFriendRequestHandler struct {
	userService       *store.UserService
	friendshipService *store.FriendshipService
	validate          *validator.Validate
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

func HandleAcceptFriendRequest(userService *store.UserService, friendshipService *store.FriendshipService, validate *validator.Validate) middlewares.HandlerWithUser {
	handler := &acceptFriendRequestHandler{
		userService:       userService,
		friendshipService: friendshipService,
		validate:          validate,
	}

	return handler.handle
}

func (h *acceptFriendRequestHandler) handle(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
	requestId, err := utils.GetIntParam(r, "requestId")
	if err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request", Cause: err}
	}

	friendship, err := h.friendshipService.GetFriendshipById(requestId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &utils.ApiError{Code: http.StatusNotFound, Message: "Friend request not found", Cause: err}
		}

		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend request", Cause: err}
	}

	if friendship.FriendID != user.ID {
		return &utils.ApiError{Code: http.StatusForbidden, Message: "You cannot accept this friend request", Cause: nil}
	}

	if friendship.Status != "pending" {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Friend request already accepted or rejected", Cause: nil}
	}

	if friendship.InviterID == user.ID {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "You cannot accept your own friend request", Cause: nil}
	}

	err = h.friendshipService.AcceptFriendRequest(requestId)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when accepting friend request", Cause: err}
	}

	return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request accepted " + strconv.Itoa(requestId)})
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
	Email string `json:"email" validate:"required,email"`
}

type AcceptFriendRequestBody struct {
	FriendID int `json:"friendId" validate:"required,number"`
}

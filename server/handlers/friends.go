package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/middlewares"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
	"time"
)

func HandleSendFriendRequest(
	userService store.UserServiceInterface,
	friendshipService store.FriendshipServiceInterface,
	validate *validator.Validate,
) middlewares.HandlerWithUser {

	return func(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
		body := &SendFriendRequestBody{}

		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
		}

		userToSendRequest, err := userService.FindUserByEmail(body.Email)

		if err != nil {
			if errors.Is(err, store.UserNotFoundError) {
				return &utils.ApiError{Code: http.StatusNotFound, Message: "User with that email not found", Cause: err}
			}
			return err
		}

		if userToSendRequest.ID == user.ID {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Cannot send friend request to yourself", Cause: nil}
		}

		existingFriendship, err := friendshipService.GetFriendshipByUsers(user.ID, userToSendRequest.ID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if existingFriendship != nil && existingFriendship.Status == "accepted" {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "You are already friends", Cause: nil}
		}

		if existingFriendship != nil && existingFriendship.Status == "pending" {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Friend request already sent", Cause: nil}
		}

		if existingFriendship != nil && existingFriendship.Status == "rejected" {
			if existingFriendship.FriendID == user.ID {
				err := friendshipService.DeleteRequestAndSendNew(
					existingFriendship.ID,
					user.ID,
					userToSendRequest.ID,
				)
				if err != nil {
					return err
				}
				return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
			}

			if !existingFriendship.StatusChangedAt.Valid {
				return &utils.ApiError{
					Code:    http.StatusInternalServerError,
					Message: "Unknown error when sending request",
					Cause:   errors.New("status changed at is null from database this should not happen if status is changed properly"),
				}
			}
			now := time.Now()
			changedAt := existingFriendship.StatusChangedAt.Time
			timeDiff := now.Sub(changedAt)
			if timeDiff < time.Hour*24*7 {
				return &utils.ApiError{Code: http.StatusBadRequest, Message: "Friend request already rejected, a week needs to pass to send another request", Cause: nil}
			}

			updateFriendshipError := friendshipService.MakeFriendshipPending(existingFriendship.ID)

			if updateFriendshipError != nil {
				return updateFriendshipError
			}

			return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
		}

		if err := friendshipService.SendFriendRequest(user.ID, userToSendRequest.ID); err != nil {
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when sending friend request", Cause: err}
		}

		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
	}
}

func HandleGetFriendRequests(friendshipService store.FriendshipServiceInterface) middlewares.HandlerWithUser {
	return func(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
		friendRequests, err := friendshipService.GetUsersFriendRequests(user.ID)

		if err != nil {
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend requests", Cause: err}
		}

		return utils.WriteJson(w, http.StatusOK, &utils.JSON{"requests": friendRequests})
	}
}

func HandleAcceptFriendRequest(friendshipService store.FriendshipServiceInterface) middlewares.HandlerWithUser {
	return func(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
		requestId, err := utils.GetIntParam(r, "requestId")
		if err != nil {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request", Cause: err}
		}

		friendship, err := friendshipService.GetFriendshipById(requestId)

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

		err = friendshipService.AcceptFriendRequest(requestId)

		if err != nil {
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when accepting friend request", Cause: err}
		}

		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request accepted"})
	}
}

func HandleRejectFriendRequest(friendshipService store.FriendshipServiceInterface) middlewares.HandlerWithUser {
	return func(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error {
		requestId, err := utils.GetIntParam(r, "requestId")
		if err != nil {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request", Cause: err}
		}
		friendshipToReject, findFriendshipError := friendshipService.GetFriendshipById(requestId)
		if findFriendshipError != nil {
			if errors.Is(findFriendshipError, sql.ErrNoRows) {
				return &utils.ApiError{Code: http.StatusNotFound, Message: "Friend request not found", Cause: findFriendshipError}
			}
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend request", Cause: findFriendshipError}
		}
		if friendshipToReject.FriendID != user.ID {
			return &utils.ApiError{Code: http.StatusForbidden, Message: "You cannot reject this friend request", Cause: nil}
		}
		if friendshipToReject.Status != "pending" {
			return &utils.ApiError{Code: http.StatusBadRequest, Message: "Friend request already accepted or rejected", Cause: nil}
		}

		rejectError := friendshipService.RejectFriendRequest(requestId)

		if rejectError != nil {
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when rejecting friend request", Cause: rejectError}
		}
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request rejected"})
	}
}

type SendFriendRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

type AcceptFriendRequestBody struct {
	FriendID int `json:"friendId" validate:"required,number"`
}

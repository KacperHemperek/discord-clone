package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"log/slog"
	"net/http"
	"time"
)

func HandleSendFriendRequest(
	userService store.UserServiceInterface,
	notificationStore store.NotificationServiceInterface,
	notificationWsService ws.NotificationServiceInterface,
	friendshipService store.FriendshipServiceInterface,
	validate *validator.Validate,
) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		body := &SendFriendRequestBody{}

		if err := utils.ReadAndValidateBody(r, body, validate); err != nil {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
		}

		userToSendRequest, err := userService.FindUserByEmail(body.Email)

		if err != nil {
			if errors.Is(err, store.UserNotFoundError) {
				return &utils.APIError{Code: http.StatusNotFound, Message: "User with that email not found", Cause: err}
			}
			return err
		}

		if userToSendRequest.ID == c.User.ID {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Cannot send friend request to yourself", Cause: nil}
		}

		existingFriendship, err := friendshipService.GetFriendshipByUsers(c.User.ID, userToSendRequest.ID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if existingFriendship != nil && existingFriendship.Status == "accepted" {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "You are already friends", Cause: nil}
		}

		if existingFriendship != nil && existingFriendship.Status == "pending" {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Friend request already sent", Cause: nil}
		}

		if existingFriendship != nil && existingFriendship.Status == "rejected" {
			if existingFriendship.FriendID == c.User.ID {
				err := friendshipService.DeleteRequestAndSendNew(
					existingFriendship.ID,
					c.User.ID,
					userToSendRequest.ID,
				)
				if err != nil {
					return err
				}
				return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
			}

			if !existingFriendship.StatusChangedAt.Valid {
				return &utils.APIError{
					Code:    http.StatusInternalServerError,
					Message: "Unknown error when sending request",
					Cause:   errors.New("status changed at is null from database this should not happen if status is changed properly"),
				}
			}
			now := time.Now()
			changedAt := existingFriendship.StatusChangedAt.Time
			timeDiff := now.Sub(changedAt)
			if timeDiff < time.Hour*24*7 {
				return &utils.APIError{Code: http.StatusBadRequest, Message: "Friend request already rejected, a week needs to pass to send another request", Cause: nil}
			}

			updateFriendshipError := friendshipService.MakeFriendshipPending(existingFriendship.ID)

			if updateFriendshipError != nil {
				return updateFriendshipError
			}
			n, err := notificationStore.CreateFriendRequestNotification(userToSendRequest.ID, models.FriendRequestNotificationData{
				TestValue: "this is a test value",
			})

			if err != nil {
				return err
			}

			sendNotificationError := notificationWsService.SendNotification(userToSendRequest.ID, n)
			if sendNotificationError != nil {
				slog.Info("could not send notification", "error", sendNotificationError)
			}

			return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
		}

		if err := friendshipService.SendFriendRequest(c.User.ID, userToSendRequest.ID); err != nil {
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when sending friend request", Cause: err}
		}

		n, err := notificationStore.CreateFriendRequestNotification(userToSendRequest.ID, models.FriendRequestNotificationData{
			TestValue: "this is a test value",
		})

		if err != nil {
			return err
		}

		sendNotificationError := notificationWsService.SendNotification(userToSendRequest.ID, n)
		if sendNotificationError != nil {
			slog.Info("could not send notification", "error", sendNotificationError)
		}

		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request sent"})
	}
}

func HandleGetFriendRequests(friendshipService store.FriendshipServiceInterface) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		friendRequests, err := friendshipService.GetUsersFriendRequests(c.User.ID)

		if err != nil {
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend requests", Cause: err}
		}

		return utils.WriteJson(w, http.StatusOK, &utils.JSON{"requests": friendRequests})
	}
}

func HandleGetFriends(friendshipService store.FriendshipServiceInterface) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		users, err := friendshipService.GetFriendsByUserID(c.User.ID)

		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &utils.JSON{"friends": users})
	}
}

func HandleAcceptFriendRequest(friendshipService store.FriendshipServiceInterface) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		requestId, err := utils.GetIntParam(r, "requestId")
		if err != nil {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Invalid request", Cause: err}
		}

		friendship, err := friendshipService.GetFriendshipByID(requestId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &utils.APIError{Code: http.StatusNotFound, Message: "Friend request not found", Cause: err}
			}

			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend request", Cause: err}
		}

		if friendship.FriendID != c.User.ID {
			return &utils.APIError{Code: http.StatusForbidden, Message: "You cannot accept this friend request", Cause: nil}
		}

		if friendship.Status != "pending" {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Friend request already accepted or rejected",
				Cause:   nil,
			}
		}

		if friendship.InviterID == c.User.ID {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "You cannot accept your own friend request",
				Cause:   nil,
			}
		}

		err = friendshipService.AcceptFriendRequest(requestId)

		if err != nil {
			return &utils.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Unknown error when accepting friend request",
				Cause:   err,
			}
		}

		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request accepted"})
	}
}

func HandleRejectFriendRequest(friendshipService store.FriendshipServiceInterface) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		requestId, err := utils.GetIntParam(r, "requestId")
		if err != nil {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Invalid request", Cause: err}
		}
		friendshipToReject, findFriendshipError := friendshipService.GetFriendshipByID(requestId)
		if findFriendshipError != nil {
			if errors.Is(findFriendshipError, sql.ErrNoRows) {
				return &utils.APIError{Code: http.StatusNotFound, Message: "Friend request not found", Cause: findFriendshipError}
			}
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when getting friend request", Cause: findFriendshipError}
		}
		if friendshipToReject.FriendID != c.User.ID {
			return &utils.APIError{Code: http.StatusForbidden, Message: "You cannot reject this friend request", Cause: nil}
		}
		if friendshipToReject.Status != "pending" {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Friend request already accepted or rejected", Cause: nil}
		}

		rejectError := friendshipService.RejectFriendRequest(requestId)

		if rejectError != nil {
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when rejecting friend request", Cause: rejectError}
		}
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"message": "Friend request rejected"})
	}
}

func HandleRemoveFriend(friendshipService store.FriendshipServiceInterface) utils.APIHandler {
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
		friendID, err := utils.GetIntParam(r, "friendID")
		if err != nil {
			return &utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "Friend id is not valid",
			}
		}
		friend, err := friendshipService.GetFriendshipByUsers(friendID, c.User.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &utils.APIError{
					Code:    http.StatusNotFound,
					Message: "User does not have a friend with that id",
				}
			}
			return err
		}
		deleteError := friendshipService.DeleteFriendship(friend.ID)
		if err != nil {
			return deleteError
		}
		return utils.WriteJson(w, http.StatusOK, &response{Message: "Friend removed successfully"})
	}
}

type SendFriendRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

type AcceptFriendRequestBody struct {
	FriendID int `json:"friendId" validate:"required,number"`
}

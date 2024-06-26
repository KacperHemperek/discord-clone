package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type RegisterUserRequest struct {
	Username        string `json:"username" validate:"required,max=24,min=2"`
	Password        string `json:"password" validate:"required,max=24,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,max=24,min=8"`
	Email           string `json:"email" validate:"required,email"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func HandleRegisterUser(userService *store.UserService, validate *validator.Validate) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, _ *utils.APIContext) error {
		body := &RegisterUserRequest{}

		if err := utils.ReadBody(r, body); err != nil {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
		}

		if err := validate.Struct(body); err != nil {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
		}

		if body.Password != body.ConfirmPassword {
			return &utils.APIError{Code: http.StatusBadRequest, Message: "Passwords do not match", Cause: nil}
		}

		existingUser, err := userService.FindUserByEmail(body.Email)

		if err != nil {
			if !errors.Is(err, store.UserNotFoundError) {
				return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when finding user", Cause: err}
			}
		}

		if existingUser != nil {
			return &utils.APIError{Code: http.StatusConflict, Message: "User with this email already exists", Cause: nil}
		}

		hashedPassword, err := utils.EncryptPassword(body.Password)

		if err != nil {
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when creating a user", Cause: err}
		}

		body.Password = hashedPassword

		user, err := userService.CreateUser(body.Username, body.Password, body.Email)

		if err != nil {
			return &utils.APIError{Code: http.StatusInternalServerError, Message: "Unknown error when creating user", Cause: err}
		}

		accessToken, accessTokenError := utils.NewAccessToken(&utils.NewTokenProps{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
		if accessTokenError != nil {
			return accessTokenError
		}
		refreshToken, refreshTokenError := utils.NewRefreshToken(&utils.NewTokenProps{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
		if refreshTokenError != nil {
			return refreshTokenError
		}

		utils.SetAuthCookies(w, accessToken, refreshToken)

		return utils.WriteJson(w, http.StatusCreated, &utils.JSON{
			"message": "user successfully registered",
			"user":    user,
		})
	}
}

func HandleLogin(userService *store.UserService, validate *validator.Validate) utils.APIHandler {

	return func(w http.ResponseWriter, r *http.Request, _ *utils.APIContext) error {
		InvalidRequestApiError := &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}

		InvalidUserOrPasswordApiError := &utils.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}

		body := &LoginUserRequest{}

		if err := utils.ReadBody(r, body); err != nil {
			return InvalidRequestApiError
		}

		if err := validate.Struct(body); err != nil {
			return InvalidRequestApiError
		}

		user, err := userService.FindUserByEmail(body.Email)

		if err != nil {
			if errors.Is(err, store.UserNotFoundError) {
				return InvalidUserOrPasswordApiError
			}

			return err
		}

		if err := utils.CheckPassword(user.Password, body.Password); err != nil {
			return InvalidUserOrPasswordApiError
		}

		accessToken, err := utils.NewAccessToken(&utils.NewTokenProps{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})

		if err != nil {
			return err
		}

		refreshToken, err := utils.NewRefreshToken(&utils.NewTokenProps{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})

		if err != nil {
			return err
		}

		utils.SetAuthCookies(w, accessToken, refreshToken)

		return utils.WriteJson(
			w,
			http.StatusOK,
			&utils.JSON{
				"message": "user successfully logged in",
				"user":    user,
			},
		)
	}
}

func HandleGetLoggedInUser() utils.APIHandler {

	type response struct {
		User         *utils.JWTUser `json:"user"`
		AccessToken  string         `json:"accessToken"`
		RefreshToken string         `json:"refreshToken"`
	}

	return func(w http.ResponseWriter, r *http.Request, c *utils.APIContext) error {
		accessToken, err := utils.GetAccessToken(r)
		if err != nil {
			return err
		}
		refreshToken, err := utils.GetRefreshToken(r)
		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, &response{
			User:         c.User,
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		})
	}
}

func HandleLogoutUser() utils.APIHandler {

	return func(w http.ResponseWriter, _ *http.Request, _ *utils.APIContext) error {
		utils.RemoveAuthTokensCookies(w)

		return utils.WriteJson(w, http.StatusOK, &utils.JSON{"message": "user successfully logged out"})
	}
}

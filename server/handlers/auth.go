package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/middlewares"
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

type registerUserHandler struct {
	userService *store.UserService
	validator   *validator.Validate
}

type loginHandler struct {
	userService *store.UserService
	validate    *validator.Validate
}

type getLoggedInUserHandler struct{}

type logoutUserHandler struct{}

func HandleRegisterUser(userService *store.UserService, validate *validator.Validate) utils.Handler {
	handler := &registerUserHandler{
		userService: userService,
		validator:   validate,
	}

	return handler.handle
}

func (h *registerUserHandler) handle(w http.ResponseWriter, r *http.Request) error {
	body := &RegisterUserRequest{}

	if err := utils.ReadBody(r, body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
	}

	if err := h.validator.Struct(body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
	}

	if body.Password != body.ConfirmPassword {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Passwords do not match", Cause: nil}
	}

	existingUser, err := h.userService.FindUserByEmail(body.Email)

	if err != nil {
		if !errors.Is(err, store.UserNotFoundError) {
			return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when finding user", Cause: err}
		}
	}

	if existingUser != nil {
		return &utils.ApiError{Code: http.StatusConflict, Message: "User with this email already exists", Cause: nil}
	}

	hashedPassword, err := utils.EncryptPassword(body.Password)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when creating a user", Cause: err}
	}

	body.Password = hashedPassword

	user, err := h.userService.CreateUser(body.Username, body.Password, body.Email)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when creating user", Cause: err}
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

func HandleLogin(userService *store.UserService, validate *validator.Validate) utils.Handler {
	handler := &loginHandler{
		userService: userService,
		validate:    validate,
	}

	return handler.handle
}

func (h *loginHandler) handle(w http.ResponseWriter, r *http.Request) error {
	InvalidRequestApiError := &utils.ApiError{
		Code:    http.StatusBadRequest,
		Message: "Invalid request body",
	}

	InvalidUserOrPasswordApiError := &utils.ApiError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid email or password",
	}

	body := &LoginUserRequest{}

	if err := utils.ReadBody(r, body); err != nil {
		return InvalidRequestApiError
	}

	if err := h.validate.Struct(body); err != nil {
		return InvalidRequestApiError
	}

	user, err := h.userService.FindUserByEmail(body.Email)

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

func HandleGetLoggedInUser() middlewares.HandlerWithUser {
	handler := &getLoggedInUserHandler{}

	return handler.handle
}

func (h *getLoggedInUserHandler) handle(w http.ResponseWriter, _ *http.Request, user *utils.JWTUser) error {
	return utils.WriteJson(w, http.StatusOK, &utils.JSON{"user": user})
}

func HandleLogoutUser() utils.Handler {
	handler := &logoutUserHandler{}

	return handler.Handle
}

func (h *logoutUserHandler) Handle(w http.ResponseWriter, _ *http.Request) error {
	utils.RemoveAuthTokensCookies(w)

	return utils.WriteJson(w, http.StatusOK, &utils.JSON{"message": "user successfully logged out"})
}

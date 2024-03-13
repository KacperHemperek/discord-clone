package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=24,min=2"`
	Password string `json:"password" validate:"required,max=24,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserHandler struct {
	userService *store.UserService
	validator   *validator.Validate
}

type LoginHandler struct {
	userService *store.UserService
	validator   *validator.Validate
}

type GetLoggedInUserHandler struct {
	userService *store.UserService
}

type LogoutUserHandler struct{}

func (h *RegisterUserHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	body := &RegisterUserRequest{}

	if err := utils.ReadBody(r, body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
	}

	if err := h.validator.Struct(body); err != nil {
		return &utils.ApiError{Code: http.StatusBadRequest, Message: "Invalid request body", Cause: err}
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

	accessToken, accessTokenError := utils.NewAccessToken(user.ID, user.Username, user.Email)
	if accessTokenError != nil {
		return accessTokenError
	}
	refreshToken, refreshTokenError := utils.NewRefreshToken(user.ID, user.Username, user.Email)
	if refreshTokenError != nil {
		return refreshTokenError
	}

	utils.SetAuthTokens(w, accessToken, refreshToken)

	return utils.WriteJson(w, http.StatusCreated, user)
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.validator.Struct(body); err != nil {
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

	accessToken, err := utils.NewAccessToken(user.ID, user.Username, user.Email)

	if err != nil {
		return err
	}

	refreshToken, err := utils.NewRefreshToken(user.ID, user.Username, user.Email)

	if err != nil {
		return err
	}

	utils.SetAuthTokens(w, accessToken, refreshToken)

	return utils.WriteJson(
		w,
		http.StatusOK,
		&utils.JSON{
			"message": "user successfully logged in",
			"user":    user,
		},
	)
}

func (h *GetLoggedInUserHandler) Handle(w http.ResponseWriter, _ *http.Request, user *models.User) error {
	return utils.WriteJson(w, http.StatusOK, &utils.JSON{"user": user})
}

func (h *LogoutUserHandler) Handle(w http.ResponseWriter, _ *http.Request) error {
	utils.RemoveAuthTokensCookies(w)

	return utils.WriteJson(w, http.StatusOK, &utils.JSON{"message": "user successfully logged out"})
}

type RegisterUserParams struct {
	UserService *store.UserService
	Validator   *validator.Validate
}

type LoginUserParams struct {
	UserService *store.UserService
	Validator   *validator.Validate
}

type GetLoggedInUserParams struct {
	UserService *store.UserService
}

func NewRegisterUserHandler(p *RegisterUserParams) *RegisterUserHandler {
	return &RegisterUserHandler{
		userService: p.UserService,
		validator:   p.Validator,
	}
}

func NewLoginHandler(p *LoginUserParams) *LoginHandler {
	return &LoginHandler{
		userService: p.UserService,
		validator:   p.Validator,
	}
}

func NewGetLoggedInUserHandler(p *GetLoggedInUserParams) *GetLoggedInUserHandler {
	return &GetLoggedInUserHandler{
		userService: p.UserService,
	}
}

func NewLogoutUserHandler() *LogoutUserHandler {
	return &LogoutUserHandler{}
}

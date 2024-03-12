package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := encryptPassword(body.Password)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when creating a user", Cause: err}
	}

	body.Password = hashedPassword

	user, err := h.userService.CreateUser(body.Username, body.Password, body.Email)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when creating user", Cause: err}
	}

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

		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when signing user in", Cause: err}
	}

	if err := checkPassword(user.Password, body.Password); err != nil {
		return InvalidUserOrPasswordApiError
	}

	return utils.WriteJson(
		w,
		http.StatusOK,
		&utils.JSON{
			"message": "user successfully logged in",
			"user":    user,
		},
	)
}

type RegisterUserParams struct {
	UserService *store.UserService
	Validator   *validator.Validate
}

type LoginUserParams struct {
	UserService *store.UserService
	Validator   *validator.Validate
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

func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

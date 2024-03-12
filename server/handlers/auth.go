package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=24,min=2"`
	Password string `json:"password" validate:"required,max=24,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type RegisterUserHandler struct {
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

	user, err := h.userService.CreateUser(body.Username, body.Password, body.Email)

	if err != nil {
		return &utils.ApiError{Code: http.StatusInternalServerError, Message: "Unknown error when creating user", Cause: err}
	}

	return utils.WriteJson(w, http.StatusCreated, user)
}

type RegisterUserParams struct {
	UserService *store.UserService
	Validator   *validator.Validate
}

func NewRegisterUserHandler(p *RegisterUserParams) *RegisterUserHandler {
	return &RegisterUserHandler{
		userService: p.UserService,
		validator:   p.Validator,
	}
}

func encryptPassword(password string) string {

	return password
}

package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/store"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type AuthMiddleware struct {
	userService *store.UserService
}

var unauthorizedApiError = &utils.ApiError{
	Code:    http.StatusUnauthorized,
	Message: "Unauthorized",
}

func (m *AuthMiddleware) Use(h HandlerWithUser) utils.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		accessToken, err := utils.GetAccessToken(r)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				return unauthorizedApiError
			}
			refreshTokenUser, accessToken, newRefreshToken, err := createNewAccessTokenAndRefreshToken(r)

			if err != nil {
				return unauthorizedApiError
			}

			utils.SetAuthTokens(w, accessToken, newRefreshToken)

			user, err := m.userService.FindUserByEmail(refreshTokenUser.Email)

			if err != nil {
				return unauthorizedApiError
			}

			return h(w, r, user)
		}

		accessTokenUser, err := utils.ParseUserToken(accessToken)

		if err != nil {
			if !errors.Is(err, jwt.ErrTokenExpired) {
				return unauthorizedApiError
			}
			refreshTokenUser, accessToken, newRefreshToken, err := createNewAccessTokenAndRefreshToken(r)

			if err != nil {
				return unauthorizedApiError
			}

			utils.SetAuthTokens(w, accessToken, newRefreshToken)

			accessTokenUser = refreshTokenUser
		}

		userFromDb, err := m.userService.FindUserByEmail(accessTokenUser.Email)

		if err != nil {
			return unauthorizedApiError
		}
		return h(w, r, userFromDb)
	}
}

func createNewAccessTokenAndRefreshToken(r *http.Request) (user *utils.JWTUser, accessToken, refreshToken string, err error) {
	oldRefreshToken, err := utils.GetRefreshToken(r)
	if err != nil {
		return nil, "", "", fmt.Errorf("error when getting refresh token")
	}

	user, err = utils.ParseUserToken(oldRefreshToken)

	if err != nil {
		return nil, "", "", fmt.Errorf("error when parsing refresh token")
	}

	accessToken, err = utils.NewAccessToken(user.UserID, user.Username, user.Email)

	if err != nil {
		return nil, "", "", fmt.Errorf("error when creating new access token")
	}

	refreshToken, err = utils.NewRefreshToken(user.UserID, user.Username, user.Email)

	if err != nil {
		return nil, "", "", fmt.Errorf("error when creating new refresh token")
	}

	return user, accessToken, refreshToken, nil
}

type AuthMiddlewareParams struct {
	UserService *store.UserService
}

func NewAuthMiddleware(params *AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		userService: params.UserService,
	}
}

type HandlerWithUser = func(w http.ResponseWriter, r *http.Request, user *models.User) error

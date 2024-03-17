package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type AuthMiddleware struct{}

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
			user, accessToken, newRefreshToken, err := createNewAccessTokenAndRefreshToken(r)

			if err != nil {
				return unauthorizedApiError
			}

			utils.SetAuthTokens(w, accessToken, newRefreshToken)

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

		return h(w, r, accessTokenUser)
	}
}

func createNewAccessTokenAndRefreshToken(r *http.Request) (user *utils.JWTUser, accessToken, refreshToken string, err error) {
	oldRefreshToken, err := utils.GetRefreshToken(r)
	if err != nil {
		return nil, "", "", err
	}

	user, err = utils.ParseUserToken(oldRefreshToken)

	if err != nil {
		return nil, "", "", fmt.Errorf("error when parsing refresh token")
	}

	token := &utils.NewTokenProps{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	accessToken, err = utils.NewAccessToken(token)

	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err = utils.NewRefreshToken(token)

	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

type HandlerWithUser = func(w http.ResponseWriter, r *http.Request, user *utils.JWTUser) error

package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

var unauthorizedApiError = &utils.APIError{
	Code:    http.StatusUnauthorized,
	Message: "Unauthorized",
}

type AuthMiddleware = func(h utils.APIHandler) utils.APIHandler

func NewAuthMiddleware() AuthMiddleware {
	return func(h utils.APIHandler) utils.APIHandler {
		return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
			accessToken, err := utils.GetAccessToken(r)
			if err != nil {
				if !errors.Is(err, http.ErrNoCookie) {
					return unauthorizedApiError
				}
				user, accessToken, newRefreshToken, err := createNewAccessTokenAndRefreshToken(r)

				if err != nil {
					return unauthorizedApiError
				}

				utils.SetAuthCookies(w, accessToken, newRefreshToken)

				if err != nil {
					return unauthorizedApiError
				}

				c.User = user

				return h(w, r, c)
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

				utils.SetAuthCookies(w, accessToken, newRefreshToken)

				accessTokenUser = refreshTokenUser
			}

			c.User = accessTokenUser

			return h(w, r, c)
		}
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

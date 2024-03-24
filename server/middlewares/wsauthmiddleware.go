package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/utils"
	"github.com/kacperhemperek/discord-go/ws"
	"net/http"
)

type WsAuthMiddleware = func(h utils.APIHandler) utils.APIHandler

func NewWsAuthMiddleware() WsAuthMiddleware {
	return func(h utils.APIHandler) utils.APIHandler {
		return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
			accessToken, err := getAccessTokenFromQueryParams(r)
			if err != nil {
				return unauthorizedApiError
			}
			var newAccessToken string
			var newRefreshToken string
			accessTokenUser, err := utils.ParseUserToken(accessToken)
			if err != nil {
				if !errors.Is(err, jwt.ErrTokenExpired) {
					return unauthorizedApiError
				}
				refreshTokenUser, accessToken, refreshToken, err := wsCreateNewAccessAndRefreshToken(r)
				if err != nil {
					return unauthorizedApiError
				}
				newAccessToken = accessToken
				newRefreshToken = refreshToken
				accessTokenUser = refreshTokenUser
				if err != nil {

					return err
				}
			}
			upgrader := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return err
			}
			c.Conn = conn
			if newAccessToken != "" && newRefreshToken != "" {
				_ = c.Conn.WriteJSON(NewUpdateTokenMessage(newAccessToken, newRefreshToken))
			}
			c.User = accessTokenUser
			return h(w, r, c)
		}
	}
}

func wsCreateNewAccessAndRefreshToken(r *http.Request) (user *utils.JWTUser, accessToken, refreshToken string, err error) {
	oldRefreshToken, err := getRefreshTokenFromQueryParams(r)
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

type UpdateTokenMessage struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewUpdateTokenMessage(at, rt string) *UpdateTokenMessage {
	return &UpdateTokenMessage{
		Type:         ws.UpdateAccessToken,
		AccessToken:  at,
		RefreshToken: rt,
	}
}

func getAccessTokenFromQueryParams(r *http.Request) (string, error) {
	accessToken := r.URL.Query().Get("accessToken")

	if accessToken == "" {
		return "", errTokenNotFound
	}

	return accessToken, nil
}

func getRefreshTokenFromQueryParams(r *http.Request) (string, error) {
	refreshToken := r.URL.Query().Get("refreshToken")

	if refreshToken == "" {
		return "", errTokenNotFound
	}

	return refreshToken, nil
}

var errTokenNotFound = errors.New("token not found in query params")

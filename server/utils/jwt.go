package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

var (
	AccessTokenCookieName  = "accessToken"
	RefreshTokenCookieName = "refreshToken"
)

type JWTUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	jwt.RegisteredClaims
}

func newAccessTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
}

func newRefreshTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
}

func SetAuthTokens(w http.ResponseWriter, accessToken string, refreshToken string) {

	tokenCookie := newAccessTokenCookie(accessToken)
	refreshCookie := newRefreshTokenCookie(refreshToken)

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, refreshCookie)
}

type NewTokenProps struct {
	ID        int
	Email     string
	Username  string
	CreatedAt string
	UpdatedAt string
}

func NewAccessToken(u *NewTokenProps) (string, error) {
	user := &JWTUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: getAccessTokenExpiryDate(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	secret, err := getJWTSecret()

	if err != nil {
		return "", err
	}

	return token.SignedString([]byte(secret))
}

func NewRefreshToken(u *NewTokenProps) (string, error) {
	user := &JWTUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: getRefreshTokenExpiryDate(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	secret, err := getJWTSecret()

	if err != nil {
		return "", err
	}

	return token.SignedString([]byte(secret))
}

func ParseUserToken(tokenString string) (*JWTUser, error) {
	token, err := jwt.Parse(tokenString, parseFunc, jwt.WithExpirationRequired())

	if err != nil {
		fmt.Println("Error when parsing token: ", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &JWTUser{}
		user.ID = int(claims["id"].(float64))
		user.Username = claims["username"].(string)
		user.Email = claims["email"].(string)
		user.CreatedAt = claims["createdAt"].(string)
		user.UpdatedAt = claims["updatedAt"].(string)

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func parseFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	secret, err := getJWTSecret()

	if err != nil {
		return nil, err
	}

	return []byte(secret), nil
}

func getJWTSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	return secret, nil
}

func getAccessTokenExpiryDate() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
}

func getRefreshTokenExpiryDate() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
}

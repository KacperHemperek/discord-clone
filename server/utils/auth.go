package utils

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func GetAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func getTokenCookieWithExpiredDate(tokenName string) *http.Cookie {
	expired := time.Now().Add(-1 * 7 * time.Hour)

	return &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Path:     "/",
		Expires:  expired,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
}

func RemoveAuthTokensCookies(w http.ResponseWriter) {
	http.SetCookie(w, getTokenCookieWithExpiredDate(AccessTokenCookieName))
	http.SetCookie(w, getTokenCookieWithExpiredDate(RefreshTokenCookieName))
}

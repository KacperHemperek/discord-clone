package middlewares

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	testUsername = "test_username"
	testEmail    = "test@email.com"
	testID       = 2137
)

func TestAuthMiddleware_BothTokensAreValid(t *testing.T) {
	accessTokenExp := time.Now().Add(time.Hour)
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7)

	rr, _, err := setupAuthMiddlewareTestWithTokens(t, accessTokenExp, refreshTokenExp)
	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	response := rr.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	jsonResponse := make(map[string]interface{})

	err = json.NewDecoder(response.Body).Decode(&jsonResponse)

	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	if id, ok := jsonResponse["id"]; !ok || id != float64(69) {
		t.Errorf("Expected id to be %d, got %v", testID, id)
	}

	if username, ok := jsonResponse["username"]; !ok || username != testUsername {
		t.Errorf("Expected testUsername to be %s, got %v", testUsername, username)
	}

	if email, ok := jsonResponse["email"]; !ok || email != testEmail {
		t.Errorf("Expected testEmaill to be %s, got %v", testEmail, email)
	}
}

func TestAuthMiddleware_AccessTokenIsInvalid(t *testing.T) {
	accessTokenExp := time.Now().Add(-time.Hour)
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7)

	rr, _, err := setupAuthMiddlewareTestWithTokens(t, accessTokenExp, refreshTokenExp)
	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	result := rr.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

	jsonResponse := make(map[string]interface{})
	err = json.NewDecoder(result.Body).Decode(&jsonResponse)

	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	if id, ok := jsonResponse["id"]; !ok || id != float64(69) {
		t.Errorf("Expected id to be %d, got %v", testID, id)
	}

	if username, ok := jsonResponse["username"]; !ok || username != testUsername {
		t.Errorf("Expected testUsername to be %s, got %v", testUsername, username)
	}

	if email, ok := jsonResponse["email"]; !ok || email != testEmail {
		t.Errorf("Expected testEmaill to be %s, got %v", testEmail, email)
	}

	cookies := result.Cookies()

	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	accessTokenCookieValue := findCookie(cookies, utils.AccessTokenCookieName)

	if accessTokenCookieValue == nil {
		t.Errorf("Expected access token cookie to be set")
	}

	refreshTokenCookieValue := findCookie(cookies, utils.RefreshTokenCookieName)

	if refreshTokenCookieValue == nil {
		t.Errorf("Expected refresh token cookie to be set")
	}

}

func TestAuthMiddleware_AccessTokenAndRefreshTokenAreExpired(t *testing.T) {
	accessTokenExp := time.Now().Add(-time.Hour)
	refreshTokenExp := time.Now().Add(-time.Hour * 24 * 7)

	rr, _, err := setupAuthMiddlewareTestWithTokens(t, accessTokenExp, refreshTokenExp)
	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	result := rr.Result()

	if result.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, result.StatusCode)
	}

	expectedMessage := "Unauthorized"
	expectedCode := http.StatusUnauthorized
	jsonResponse := make(map[string]interface{})
	err = json.NewDecoder(result.Body).Decode(&jsonResponse)

	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	if message, ok := jsonResponse["message"]; !ok || message != expectedMessage {
		t.Errorf("Expected message to be %s, got %v", expectedMessage, message)
	}

	if code, ok := jsonResponse["code"]; !ok || code != float64(expectedCode) {
		t.Errorf("Expected code to be %d, got %v", expectedCode, code)
	}
}

func TestAuthMiddleware_AccessTokenMissingAndRefreshTokenValid(t *testing.T) {
	rr, err := setupAuthMiddlewareTestWithoutAccessToken(t, time.Now().Add(time.Hour*24*7))

	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	result := rr.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

	cookies := result.Cookies()

	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	accessTokenCookieValue := findCookie(cookies, utils.AccessTokenCookieName)

	if accessTokenCookieValue == nil {
		t.Errorf("Expected access token cookie to be set")
	}

	refreshTokenCookieValue := findCookie(cookies, utils.RefreshTokenCookieName)

	if refreshTokenCookieValue == nil {
		t.Errorf("Expected refresh token cookie to be set")
	}
}

func TestAuthMiddleware_AccessTokenMissingAndRefreshTokenExpired(t *testing.T) {
	rr, err := setupAuthMiddlewareTestWithoutAccessToken(t, time.Now().Add(time.Hour*24*7))

	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	result := rr.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

	cookies := result.Cookies()

	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	accessTokenCookieValue := findCookie(cookies, utils.AccessTokenCookieName)

	if accessTokenCookieValue == nil {
		t.Errorf("Expected access token cookie to be set")
	}

	refreshTokenCookieValue := findCookie(cookies, utils.RefreshTokenCookieName)

	if refreshTokenCookieValue == nil {
		t.Errorf("Expected refresh token cookie to be set")
	}
}

func TestAuthMiddleware_RefreshTokenAndAccessTokenAreMissing(t *testing.T) {
	rr, err := setupAuthMiddlewareTestWithoutTokens(t)

	if err != nil {
		t.Errorf("Error setting up test: %s", err)
	}

	result := rr.Result()

	if result.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, result.StatusCode)
	}

	expectedMessage := "Unauthorized"
	expectedCode := http.StatusUnauthorized
	jsonResponse := make(map[string]interface{})
	err = json.NewDecoder(result.Body).Decode(&jsonResponse)

	if err != nil {
		t.Errorf("Error decoding response body: %s", err)
	}

	if message, ok := jsonResponse["message"]; !ok || message != expectedMessage {
		t.Errorf("Expected message to be %s, got %v", expectedMessage, message)
	}

	if code, ok := jsonResponse["code"]; !ok || code != float64(expectedCode) {
		t.Errorf("Expected code to be %d, got %v", expectedCode, code)
	}

	accessTokenCookie := findCookie(result.Cookies(), utils.AccessTokenCookieName)

	if accessTokenCookie != nil {
		t.Errorf("Expected access token cookie to not be set")
	}

	refreshTokenCookie := findCookie(result.Cookies(), utils.RefreshTokenCookieName)

	if refreshTokenCookie != nil {
		t.Errorf("Expected refresh token cookie to not be set")
	}
}

func testHandler(w http.ResponseWriter, _ *http.Request, c *utils.Context) error {
	return utils.WriteJson(w, http.StatusOK, c.User)
}

func setupAuthMiddlewareTestWithTokens(t *testing.T, accessTokenExp, refreshTokenExp time.Time) (rr *httptest.ResponseRecorder, req *http.Request, err error) {
	t.Setenv("JWT_SECRET", "test_secret")
	authMiddleware := NewAuthMiddleware()
	rr = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		return nil, nil, err
	}

	accessTokenUser := getJWTUser(accessTokenExp)

	refreshTokenUser := getJWTUser(refreshTokenExp)

	accessToken, err := getTestToken(accessTokenUser)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := getTestToken(refreshTokenUser)

	if err != nil {
		return nil, nil, err
	}

	req.AddCookie(getTokenCookie(accessToken, utils.AccessTokenCookieName))
	req.AddCookie(getTokenCookie(refreshToken, utils.RefreshTokenCookieName))

	handler := utils.HandlerFunc(authMiddleware(testHandler))

	handler.ServeHTTP(rr, req)

	return rr, req, nil
}

func setupAuthMiddlewareTestWithoutAccessToken(t *testing.T, refreshTokenExp time.Time) (rr *httptest.ResponseRecorder, err error) {
	t.Setenv("JWT_SECRET", "test_secret")
	authMiddleware := NewAuthMiddleware()
	rr = httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		return nil, err
	}

	refreshTokenUser := getJWTUser(refreshTokenExp)

	refreshToken, err := getTestToken(refreshTokenUser)

	if err != nil {
		return nil, err
	}

	req.AddCookie(getTokenCookie(refreshToken, utils.RefreshTokenCookieName))

	handler := utils.HandlerFunc(authMiddleware(testHandler))

	handler.ServeHTTP(rr, req)

	return rr, nil
}

func setupAuthMiddlewareTestWithoutTokens(t *testing.T) (rr *httptest.ResponseRecorder, err error) {
	t.Setenv("JWT_SECRET", "test_secret")
	authMiddleware := NewAuthMiddleware()
	rr = httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		return nil, err
	}

	handler := utils.HandlerFunc(authMiddleware(testHandler))

	handler.ServeHTTP(rr, req)

	return rr, nil
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

func getTestToken(user *utils.JWTUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func getTokenCookie(token, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
}

func getJWTUser(expiresAt time.Time) *utils.JWTUser {
	return &utils.JWTUser{
		ID:        69,
		Username:  testUsername,
		Email:     testEmail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}

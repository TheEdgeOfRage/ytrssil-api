package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	db_mock "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/mocks/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *http.Server {
	db := &db_mock.DBMock{
		AuthenticateUserFunc: func(ctx context.Context, username, password string) (bool, error) {
			return username == "username" && password == "password", nil
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Middlewares are executed top to bottom in a stack-like manner
	router.Use(
		gin.Recovery(), // Recovery needs to go before other middlewares to catch panics
		AuthMiddleware(db),
	)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return &http.Server{Handler: router}
}

func TestSuccessfulAuthentication(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=") // username:password
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestMissingAuthorizationHeader(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"invalid authorization header"}`, w.Body.String())
}

func TestWrongCredentials(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Basic d3Jvbmc=")
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"invalid API Key"}`, w.Body.String())
}

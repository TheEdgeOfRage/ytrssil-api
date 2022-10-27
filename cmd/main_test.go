package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

var testConfig config.Config

func init() {
	testConfig = config.TestConfig()
}

func setupTestServer(t *testing.T, authEnabled bool) *http.Server {
	l := log.NewNopLogger()

	db, err := db.NewPSQLDB(l, testConfig.DB)
	if !assert.NoError(t, err) {
		return nil
	}

	handler := handler.New(l, db)
	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(
		l,
		handler,
		auth.AuthMiddleware(db),
	)
	if !assert.NoError(t, err) {
		return nil
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", testConfig.Gin.Port),
		Handler: router,
	}
}

func TestHealthz(t *testing.T) {
	server := setupTestServer(t, false)
	if !assert.NotNil(t, server) {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "healthy", w.Body.String())
}

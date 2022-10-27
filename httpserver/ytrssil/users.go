package ytrssil

import (
	"net/http"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
	"github.com/gin-gonic/gin"
)

func (s *server) CreateUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.handler.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "user created"})
}

package ytrssil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

func (s *server) SubscribeToChannel(c *gin.Context) {
	var channel models.Channel
	err := c.ShouldBindUri(&channel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.GetString("username")

	err = s.handler.SubscribeToChannel(c.Request.Context(), username, channel.ID)
	if err != nil {
		if errors.Is(err, db.ErrAlreadySubscribed) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, feedparser.ErrInvalidChannelID) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "subscribed to channel successfully"})
}

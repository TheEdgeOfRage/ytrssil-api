package ytrssil

import (
	"net/http"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
	"github.com/gin-gonic/gin"
)

func (s *server) GetNewVideos(c *gin.Context) {
	username := c.GetString("username")
	videos, err := s.handler.GetNewVideos(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.VideosResponse{
		Videos: videos,
	})
}

func (s *server) GetWatchedVideos(c *gin.Context) {
	username := c.GetString("username")
	videos, err := s.handler.GetWatchedVideos(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.VideosResponse{
		Videos: videos,
	})
}

func (s *server) FetchVideos(c *gin.Context) {
	err := s.handler.FetchVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "videos fetched successfully"})
}

func (s *server) MarkVideoAsWatched(c *gin.Context) {
	username := c.GetString("username")
	var req models.VideoURIRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.handler.MarkVideoAsWatched(c.Request.Context(), username, req.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "marked video as watched"})
}

func (s *server) MarkVideoAsUnwatched(c *gin.Context) {
	username := c.GetString("username")
	var req models.VideoURIRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.handler.MarkVideoAsUnwatched(c.Request.Context(), username, req.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "cleared video from watch history"})
}

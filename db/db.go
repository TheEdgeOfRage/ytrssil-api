package db

import (
	"context"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

// DB represents a database layer for getting video and channel data
type DB interface {
	// GetNewVideos returns unwatched videos from all channels
	GetNewVideos(ctx context.Context, username string) ([]*models.Video, error)
	AuthenticateUser(ctx context.Context, username string, password string) (bool, error)
	CreateUser(ctx context.Context, user models.User) error
}

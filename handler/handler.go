package handler

import (
	"context"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

type Handler interface {
	CreateUser(ctx context.Context, user models.User) error
	SubscribeToChannel(ctx context.Context, username string, channelID string) error
	GetNewVideos(ctx context.Context, username string) ([]models.Video, error)
	GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error)
	FetchVideos(ctx context.Context) error
	MarkVideoAsWatched(ctx context.Context, username string, videoID string) error
	MarkVideoAsUnwatched(ctx context.Context, username string, videoID string) error
}

type handler struct {
	log log.Logger
	db  db.DB
}

func New(log log.Logger, db db.DB) *handler {
	return &handler{log: log, db: db}
}

package handler

import (
	"context"

	"github.com/alexedwards/argon2id"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

type Handler interface {
	CreateUser(ctx context.Context, user models.User) error
	GetNewVideos(ctx context.Context, username string) ([]*models.Video, error)
}

type handler struct {
	log log.Logger
	db  db.DB
}

func New(log log.Logger, db db.DB) *handler {
	return &handler{log: log, db: db}
}

func (h *handler) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return h.db.CreateUser(ctx, user)
}

func (h *handler) GetNewVideos(ctx context.Context, username string) ([]*models.Video, error) {
	return h.db.GetNewVideos(ctx, username)
}

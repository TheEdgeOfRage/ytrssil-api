package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/alexedwards/argon2id"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

type Handler interface {
	CreateUser(ctx context.Context, user models.User) error
	SubscribeToChannel(ctx context.Context, username string, channelID string) error
	GetNewVideos(ctx context.Context, username string) ([]models.Video, error)
	FetchVideos(ctx context.Context) error
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

func (h *handler) SubscribeToChannel(ctx context.Context, username string, channelID string) error {
	parsedChannel, err := feedparser.Parse(h.log, channelID)
	if err != nil {
		return err
	}

	channel := models.Channel{
		ID:   channelID,
		Name: parsedChannel.Name,
	}

	err = h.db.CreateChannel(ctx, channel)
	if !errors.Is(err, db.ErrChannelExists) {
		return err
	}

	return h.db.SubscribeUserToChannel(ctx, username, channelID)
}

func (h *handler) GetNewVideos(ctx context.Context, username string) ([]models.Video, error) {
	return h.db.GetNewVideos(ctx, username)
}

func (h *handler) addVideoToAllSubscribers(ctx context.Context, channelID string, videoID string) error {
	subs, err := h.db.GetChannelSubscribers(ctx, channelID)
	if err != nil {
		h.log.Log("level", "ERROR", "call", "db.GetChannelSubscribers", "err", err)
		return err
	}

	for _, sub := range subs {
		err = h.db.AddVideoToUser(ctx, sub, videoID)
		if err != nil {
			h.log.Log("level", "ERROR", "call", "db.AddVideoToUser", "err", err)
			continue
		}
	}

	return nil
}

func (h *handler) fetchVideosForChannel(ctx context.Context, channelID string, parsedChannel *feedparser.Channel) {
	for _, parsedVideo := range parsedChannel.Videos {
		date, err := parsedVideo.Published.Parse()
		if err != nil {
			h.log.Log("level", "WARNING", "call", "feedparser.Parse", "err", err)
			continue
		}

		id := strings.Split(parsedVideo.ID, ":")[2]
		video := models.Video{
			ID:            id,
			Title:         parsedVideo.Title,
			PublishedTime: date,
		}
		err = h.db.CreateVideo(ctx, video, channelID)
		if err != nil {
			if !errors.Is(err, db.ErrVideoExists) {
				h.log.Log("level", "WARNING", "call", "db.CreateVideo", "err", err)
			}
			continue
		}
		err = h.addVideoToAllSubscribers(ctx, channelID, id)
		if err != nil {
			continue
		}
	}
}

func (h *handler) FetchVideos(ctx context.Context) error {
	channels, err := h.db.ListChannels(ctx)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		parsedChannel, err := feedparser.Parse(h.log, channel.ID)
		if err != nil {
			continue
		}

		h.fetchVideosForChannel(ctx, channel.ID, parsedChannel)
	}

	return nil
}

package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) GetNewVideos(ctx context.Context, username string) ([]models.Video, error) {
	return h.db.GetNewVideos(ctx, username)
}

func (h *handler) GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error) {
	return h.db.GetWatchedVideos(ctx, username)
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
		err = h.db.AddVideo(ctx, video, channelID)
		if err != nil {
			if !errors.Is(err, db.ErrVideoExists) {
				h.log.Log("level", "WARNING", "call", "db.AddVideo", "err", err)
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

func (h *handler) MarkVideoAsWatched(ctx context.Context, username string, videoID string) error {
	watchTime := time.Now()
	return h.db.SetVideoWatchTime(ctx, username, videoID, &watchTime)
}

func (h *handler) MarkVideoAsUnwatched(ctx context.Context, username string, videoID string) error {
	return h.db.SetVideoWatchTime(ctx, username, videoID, nil)
}

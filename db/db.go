package db

import (
	"context"
	"errors"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var (
	ErrChannelExists     = errors.New("channel already exists")
	ErrAlreadySubscribed = errors.New("already subscribed to channel")
	ErrVideoExists       = errors.New("video already exists")
)

// DB represents a database layer for getting video and channel data
type DB interface {
	// GetNewVideos returns unwatched videos from all channels
	GetNewVideos(ctx context.Context, username string) ([]models.Video, error)
	// CreateVideo adds a newly published video to the database
	CreateVideo(ctx context.Context, video models.Video, channelID string) error

	// CreateChannel starts tracking a new channel and fetch new videos for it
	CreateChannel(ctx context.Context, channel models.Channel) error
	// ListChannels lists all channels from the database
	ListChannels(ctx context.Context) ([]models.Channel, error)
	// GetChannelSubscribers lists all channels from the database
	GetChannelSubscribers(ctx context.Context, channelID string) ([]string, error)

	// AuthenticateUser verifies a user's password against a hashed value
	AuthenticateUser(ctx context.Context, user models.User) (bool, error)
	// CreateUser registers a new user in the database
	CreateUser(ctx context.Context, user models.User) error
	// DeleteUser registers a new user in the database
	DeleteUser(ctx context.Context, username string) error
	// SubscribeUserToChannel will start showing new videos for that channel to the user
	SubscribeUserToChannel(ctx context.Context, username string, channelID string) error
	// AddVideoToUser will list the video in the users feed
	AddVideoToUser(ctx context.Context, username string, videoID string) error
	// WatchVideo marks a video as watched so it no longer shows in the feed
	WatchVideo(ctx context.Context, username string, videoID string) error
}

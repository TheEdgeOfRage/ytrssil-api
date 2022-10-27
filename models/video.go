package models

import (
	"time"
)

type Video struct {
	// YouTube ID of the video
	VideoID string `json:"video_id" db:"video_id"`
	// Name of the channel the video belongs to
	ChannelName string `json:"channel_name" db:"channel_name"`
	// Title of the video
	Title string `json:"title" db:"title"`
	// Video publish timestamp
	PublishedTime time.Time `json:"published_timestamp" db:"published_timestamp"`
	// Video watch timestamp
	WatchTime *time.Time `json:"watch_timestamp" db:"watch_timestamp"`
}

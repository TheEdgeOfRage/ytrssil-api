package models

import (
	"time"
)

type Video struct {
	// YouTube ID of the video
	ID string `json:"video_id"`
	// Name of the channel the video belongs to
	ChannelName string `json:"channel_name"`
	// Title of the video
	Title string `json:"title"`
	// Video publish timestamp
	PublishedTime time.Time `json:"published_timestamp"`
	// Video watch timestamp
	WatchTime *time.Time `json:"watch_timestamp"`
}

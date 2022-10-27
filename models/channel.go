package models

type Channel struct {
	// YouTube ID of the channel
	ChannelID string `json:"channel_id" dynamodbav:"channel_id"`
	// Name of the channel
	Name string `json:"name" dynamodbav:"name"`
	// Feed is the URL for the RSS feed
	FeedURL string `json:"feed_url" dynamodbav:"feed_url"`
}

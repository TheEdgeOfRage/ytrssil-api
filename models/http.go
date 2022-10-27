package models

type GetNewVideosResponse struct {
	Videos []*Video `json:"videos"`
}

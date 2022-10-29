package db

import (
	"context"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
	"github.com/lib/pq"
)

var getNewVideosQuery = `
	SELECT
		video_id
		, title
		, published_timestamp
		, watch_timestamp
		, name as channel_name
	FROM user_videos
	LEFT JOIN videos ON video_id=videos.id
	LEFT JOIN channels ON channel_id=channels.id
	WHERE
		1=1
		AND watch_timestamp IS NULL
		AND username=$1
	ORDER BY published_timestamp
`

func (d *postgresDB) GetNewVideos(ctx context.Context, username string) ([]models.Video, error) {
	rows, err := d.db.QueryContext(ctx, getNewVideosQuery, username)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.GetNewVideos", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	videos := make([]models.Video, 0)
	for rows.Next() {
		var video models.Video
		err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.PublishedTime,
			&video.WatchTime,
			&video.ChannelName,
		)
		if err != nil {
			d.l.Log("level", "ERROR", "function", "db.GetNewVideos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

var createVideoQuery = `INSERT INTO videos (id, title, published_timestamp, channel_id) VALUES ($1, $2, $3, $4)`

func (d *postgresDB) CreateVideo(ctx context.Context, video models.Video, channelID string) error {
	_, err := d.db.ExecContext(ctx, createVideoQuery, video.ID, video.Title, video.PublishedTime, channelID)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrVideoExists
			}
		}
		d.l.Log("level", "ERROR", "function", "db.CreateVideo", "call", "sql.Exec", "error", err)
		return err
	}

	return nil
}

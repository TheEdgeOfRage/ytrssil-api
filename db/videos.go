package db

import (
	"context"

	"github.com/georgysavva/scany/v2/sqlscan"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
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

func (d *psqlDB) GetNewVideos(ctx context.Context, username string) ([]*models.Video, error) {
	var videos []*models.Video
	err := sqlscan.Select(ctx, d.db, &videos, getNewVideosQuery, username)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.GetNewVideos", "msg", "failed to query new videos", "error", err)
		return nil, err
	}

	return videos, nil
}

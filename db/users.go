package db

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/lib/pq"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var authenticateUserQuery = `SELECT password FROM users WHERE username = $1`

func (d *postgresDB) AuthenticateUser(ctx context.Context, user models.User) (bool, error) {
	row := d.db.QueryRowContext(ctx, authenticateUserQuery, user.Username)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "error", err)
		return false, err
	}

	match, err := argon2id.ComparePasswordAndHash(user.Password, hashedPassword)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AuthenticateUser", "error", err)
		return false, err
	}

	return match, nil
}

var createUserQuery = `INSERT INTO users (username, password) VALUES ($1, $2)`

func (d *postgresDB) CreateUser(ctx context.Context, user models.User) error {
	_, err := d.db.ExecContext(ctx, createUserQuery, user.Username, user.Password)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.CreateUser", "error", err)
		return err
	}

	return nil
}

var deleteUserQuery = `DELETE FROM users WHERE username = $1`

func (d *postgresDB) DeleteUser(ctx context.Context, username string) error {
	_, err := d.db.ExecContext(ctx, deleteUserQuery, username)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.DeleteUser", "error", err)
		return err
	}

	return nil
}

var subscribeUserToChannelQuery = `INSERT INTO user_subscriptions (username, channel_id) VALUES ($1, $2)`

func (d *postgresDB) SubscribeUserToChannel(ctx context.Context, username string, channelID string) error {
	_, err := d.db.ExecContext(ctx, subscribeUserToChannelQuery, username, channelID)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrAlreadySubscribed
			}
		}
		d.l.Log("level", "ERROR", "function", "db.SubscribeUserToChannel", "error", err)
		return err
	}

	return nil
}

var addVideoToUserQuery = `INSERT INTO user_videos (username, video_id) VALUES ($1, $2)`

func (d *postgresDB) AddVideoToUser(ctx context.Context, username string, videoID string) error {
	_, err := d.db.ExecContext(ctx, addVideoToUserQuery, username, videoID)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.AddVideoToUser", "error", err)
		return err
	}

	return nil
}

var watchVideoQuery = `UPDATE user_videos SET watch_timestamp = NOW() WHERE username = $2 AND video_id = $3`

func (d *postgresDB) WatchVideo(ctx context.Context, username string, videoID string) error {
	_, err := d.db.ExecContext(ctx, watchVideoQuery, username, videoID)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.WatchVideo", "error", err)
		return err
	}

	return nil
}

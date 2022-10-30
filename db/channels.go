package db

import (
	"context"

	"github.com/lib/pq"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var createChannelQuery = `INSERT INTO channels (id, name) VALUES ($1, $2)`

func (d *postgresDB) CreateChannel(ctx context.Context, channel models.Channel) error {
	_, err := d.db.ExecContext(ctx, createChannelQuery, channel.ID, channel.Name)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrChannelExists
			}
		}

		d.l.Log("level", "ERROR", "function", "db.CreateChannel", "error", err)
		return err
	}

	return nil
}

var listChannelsQuery = `SELECT id, name FROM channels`

func (d *postgresDB) ListChannels(ctx context.Context) ([]models.Channel, error) {
	rows, err := d.db.QueryContext(ctx, listChannelsQuery)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.ListChannels", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0)
	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			d.l.Log("level", "ERROR", "function", "db.ListChannels", "call", "sql.Scan", "error", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

var getChannelSubscribersQuery = `SELECT username FROM user_subscriptions WHERE channel_id = $1`

func (d *postgresDB) GetChannelSubscribers(ctx context.Context, channelID string) ([]string, error) {
	rows, err := d.db.QueryContext(ctx, getChannelSubscribersQuery, channelID)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.GetChannelSubscribers", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	subs := make([]string, 0)
	for rows.Next() {
		var sub string
		err = rows.Scan(&sub)
		if err != nil {
			d.l.Log("level", "ERROR", "function", "db.GetChannelSubscribers", "call", "sql.Scan", "error", err)
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
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

var unsubscribeUserFromChannelQuery = `DELETE FROM user_subscriptions WHERE username = $1 AND channel_id = $2`

func (d *postgresDB) UnsubscribeUserFromChannel(ctx context.Context, username string, channelID string) error {
	resp, err := d.db.ExecContext(ctx, unsubscribeUserFromChannelQuery, username, channelID)
	if err != nil {
		d.l.Log("level", "ERROR", "function", "db.SubscribeUserToChannel", "error", err)
		return err
	}

	if affected, err := resp.RowsAffected(); err != nil || affected != 1 {
		return ErrChannelNotFound
	}

	return nil
}

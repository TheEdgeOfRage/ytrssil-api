CREATE TABLE IF NOT EXISTS channels (
	id text NOT NULL PRIMARY KEY
	, name text NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
	id text NOT NULL PRIMARY KEY
	, title text NOT NULL
	, published_timestamp timestamp with time zone NOT NULL
	, channel_id text NOT NULL REFERENCES channels(id)
);

CREATE TABLE IF NOT EXISTS users (
	username text NOT NULL PRIMARY KEY
	, password text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_videos (
	username text NOT NULL REFERENCES users(username)
	, video_id text NOT NULL REFERENCES videos(id)
	, watch_timestamp timestamp with time zone
	, CONSTRAINT user_videos_pkey PRIMARY KEY (username, video_id)
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
	username text NOT NULL REFERENCES users(username)
	, channel_id text NOT NULL REFERENCES channels(id)
	, CONSTRAINT user_subscriptions_pkey PRIMARY KEY (channel_id, username)
);

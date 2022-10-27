CREATE TABLE IF NOT EXISTS channels (
	id text NOT NULL PRIMARY KEY,
	name text NOT NULL,
	feed_url text NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
	id text NOT NULL PRIMARY KEY,
	title text NOT NULL,
	published_timestamp timestamp with time zone NOT NULL,
	channel_id text NOT NULL REFERENCES channels(id)
);

CREATE TABLE IF NOT EXISTS users (
	username text NOT NULL PRIMARY KEY,
	password text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_videos (
	id SERIAL PRIMARY KEY,
	username text NOT NULL REFERENCES users(username),
	video_id text NOT NULL REFERENCES videos(id),
	watch_timestamp timestamp with time zone
);

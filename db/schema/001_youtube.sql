-- +goose Up

-- enable extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- videos table
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    published_on TIMESTAMP NOT NULL,
    thumbnail_url TEXT NOT NULL,
    provider TEXT NOT NULL, -- provides details about video provider 'youtube', 'vimeo', 'facebook' etc
    video_id TEXT UNIQUE NOT NULL, -- contains information about the actual url of the video from where we can stream it
    view_count BIGINT NOT NULL,
    like_count BIGINT NOT NULL,
    favorite_count BIGINT NOT NULL,
    comment_count BIGINT NOT NULL
);


-- +goose Down
DROP TABLE videos;

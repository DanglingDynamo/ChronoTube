-- +goose Up

-- videos table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tite TEXT,
    desc TEXT,
    published_on TIMESTAMP,
    thumbnail_url TEXT,
    provider TEXT, -- provides details about video provider 'youtube', 'vimeo', 'facebook' etc
);

-- youtube specific data goes here
CREATE TABLE IF NOT EXISTS youtube_videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video UUID NOT NULL,
    video_id TEXT, -- contains information about the actual url of the video from where we can stream it
    view_count INT,
    like_count INT,
    favorite_count INT,
    comment_count INT,
    FOREIGN KEY (video) REFERENCES videos(id) ON DELETE CASCADE;
);

-- +goose Down
DROP TABLE youtube_videos;
DROP TABLE videos;

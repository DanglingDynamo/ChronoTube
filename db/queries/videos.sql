-- name: InsertVideo :one
INSERT INTO videos(title, desc, published_on, thumbnail_url, provider) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: InsertYoutubeDetails :one
INSERT INTO youtube_videos(video, video_id, view_count, like_count, favorite_count, comment_count) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

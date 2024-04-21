-- name: InsertVideo :one
INSERT INTO videos(title, description, published_on, thumbnail_url, provider, video_id, view_count, like_count, favorite_count, comment_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;



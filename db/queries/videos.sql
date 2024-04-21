-- name: InsertVideo :one
INSERT INTO videos(title, description, published_on, thumbnail_url, provider, video_id, view_count, like_count, favorite_count, comment_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: FetchVideosPaginated :many
SELECT * FROM videos ORDER BY published_on DESC LIMIT 10 OFFSET $1; -- implement pagination (not the most optimal method)

-- TODO: Implement cursor based pagination as an optimization

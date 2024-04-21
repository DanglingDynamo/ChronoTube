-- +goose Up
CREATE INDEX IF NOT EXISTS idx_published_on_desc ON videos(published_on DESC);
CREATE INDEX IF NOT EXISTS idx_view_count ON videos(view_count);
CREATE INDEX IF NOT EXISTS idx_like_count ON videos(like_count);
CREATE INDEX IF NOT EXISTS idx_comment_count ON videos(comment_count);

-- +goose Down
DROP INDEX IF EXISTS idx_published_on_desc;
DROP INDEX IF EXISTS idx_view_count;
DROP INDEX IF EXISTS idx_comment_count;
DROP INDEX IF EXISTS idx_like_count;

package repository

import (
	"context"
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
)

type VideoRepository interface {
	FetchVideosFromAPI(query string, publishedAfter time.Time) ([]*models.Video, error)
	StoreVideos(ctx context.Context, videos []*models.Video) error
	FetchPaginatedVideos(ctx context.Context, nextPage int) ([]models.Video, error)
}

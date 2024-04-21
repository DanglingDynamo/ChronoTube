package services

import (
	"context"
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
)

type VideoService interface {
	FetchVideosFromAPI(query string, publishedAfter time.Time) ([]*models.Video, error)
	StoreVideos(ctx context.Context, videos []*models.Video) error
}

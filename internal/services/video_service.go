package services

import (
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
)

type VideoService interface {
	FetchVideos(query string, publishedAfter time.Time) ([]*models.Video, error)
}

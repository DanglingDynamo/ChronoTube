package cronjobs

import (
	"context"
	"log/slog"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/repository"
)

func StoreVideos(
	in <-chan []*models.Video,
	service repository.VideoRepository,
	errChan chan<- error,
) {
	for videos := range in {
		err := service.StoreVideos(context.Background(), videos)
		if err != nil {
			errChan <- err
		}
	}
	slog.Info("stopping store videos")
}

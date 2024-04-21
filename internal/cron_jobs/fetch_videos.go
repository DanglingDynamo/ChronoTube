package cronjobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/services"
)

func FetchVideos(
	ctx context.Context,
	duration time.Duration,
	service services.VideoService,
	queryString string,
	out chan<- []*models.Video,
	errChan chan<- error,
) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping video cron")
			ticker.Stop()
			return
		case tick := <-ticker.C:
			videos, err := service.FetchVideos(queryString, tick)
			if err != nil {
				slog.Error(err.Error())
				errChan <- err
				continue
			}

			out <- videos
		}
	}
}

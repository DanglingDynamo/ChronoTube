package cronjobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/repository"
)

func FetchVideos(
	ctx context.Context,
	duration time.Duration,
	service repository.VideoRepository,
	queryString string,
	out chan<- []*models.Video,
) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping video cron")
			ticker.Stop()
			return
		case tick := <-ticker.C:
			videos, err := service.FetchVideosFromAPI(queryString, tick)
			if err != nil {
				slog.Error(err.Error())
				if ctx.Err() != nil {
					slog.Error(context.Canceled.Error())
				}
				continue
			}

			select {
			case out <- videos:
			case <-ctx.Done():
				slog.Info("context cancelled while sending")
				return
			default:
				slog.Info("out channel closed")
			}

		}
	}
}

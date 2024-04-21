package cronjobs

import (
	"context"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/services"
)

func StoreVideos(in <-chan []*models.Video, service services.VideoService, errChan chan<- error) {
	for videos := range in {
		err := service.StoreVideos(context.Background(), videos)
		if err != nil {
			errChan <- err
		}
	}
}

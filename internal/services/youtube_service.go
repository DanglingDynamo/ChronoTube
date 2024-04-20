package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/DanglingDynamo/chronotube/pkg/youtube"
)

// Implements the VideoService interface
type YoutubeService struct {
	client *youtube.YoutubeClient
}

func NewYoutubeService(apiKey string) (*YoutubeService, error) {
	if apiKey != "" {
		return nil, errors.New("please provide an API key")
	}

	ytClient, err := youtube.NewYoutubeClient(context.Background(), os.Getenv("YOUTUBE_API_KEY"))
	if err != nil {
		return nil, err
	}

	return &YoutubeService{
		client: ytClient,
	}, nil
}

func (service *YoutubeService) FetchVideos(query string, publishedAfter time.Time) {
	err := service.client.FetchVideos(query, publishedAfter)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

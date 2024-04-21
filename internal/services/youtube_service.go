package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/DanglingDynamo/chronotube/internal/models"
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

func (service *YoutubeService) FetchVideos(
	query string,
	publishedAfter time.Time,
) ([]*models.Video, error) {
	slog.Info("Fetching Videos")
	ytVideos, err := service.client.FetchVideos(query, publishedAfter.AddDate(0, 0, -2))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	videos := make([]*models.Video, len(ytVideos))
	ids := make([]string, len(ytVideos))
	for i, video := range ytVideos {
		ids[i] = video.Id.VideoId
	}

	videoStatistics, err := service.client.FetchVideoStatistics(ids...)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for i, video := range ytVideos {
		publishedAt, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
		if err != nil { // TODO:  complete error handling
			slog.Error(err.Error())
			continue
		}

		videos[i] = &models.Video{
			Title:         video.Snippet.Title,
			Description:   video.Snippet.Description,
			PublishedOn:   publishedAt,
			ThumbnailURL:  video.Snippet.Thumbnails.Default.Url,
			Provider:      models.VideoProviderYoutube,
			VideoID:       video.Id.VideoId,
			ViewCount:     int(videoStatistics[video.Id.VideoId].ViewCount),
			LikeCount:     int(videoStatistics[video.Id.VideoId].LikeCount),
			FavoriteCount: int(videoStatistics[video.Id.VideoId].FavoriteCount),
			CommentCount:  int(videoStatistics[video.Id.VideoId].CommentCount),
		}
	}

	return videos, nil
}

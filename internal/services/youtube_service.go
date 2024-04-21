package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/DanglingDynamo/chronotube/internal/database"
	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/pkg/youtube"
)

// Implements the VideoService interface
type YoutubeService struct {
	client *youtube.YoutubeClient
	db     *database.Queries
}

func NewYoutubeService(apiKey string, db *database.Queries) (*YoutubeService, error) {
	if apiKey != "" {
		return nil, errors.New("please provide an API key")
	}

	ytClient, err := youtube.NewYoutubeClient(context.Background(), os.Getenv("YOUTUBE_API_KEY"))
	if err != nil {
		return nil, err
	}

	return &YoutubeService{
		client: ytClient,
		db:     db,
	}, nil
}

func (service *YoutubeService) FetchVideosFromAPI(
	query string,
	publishedAfter time.Time,
) ([]*models.Video, error) {
	slog.Info("Fetching Videos")
	ytVideos, err := service.client.FetchVideos(
		query,
		publishedAfter.AddDate(0, 0, -9),
	) // Search for videos that were uploaded after current time - 9 days (added 9 days so that I get some data)
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
		if err != nil {
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

// Storing Videos is specific to the service hence implemented here
func (service *YoutubeService) StoreVideos(ctx context.Context, videos []*models.Video) error {
	var pgErr *pgconn.PgError
	storeCount := 0
	for i := range videos {
		video := videos[i]
		_, err := service.db.InsertVideo(ctx, database.InsertVideoParams{
			Title:         video.Title,
			Description:   video.Description,
			PublishedOn:   video.PublishedOn,
			ThumbnailUrl:  video.ThumbnailURL,
			Provider:      string(video.Provider),
			VideoID:       video.VideoID,
			ViewCount:     int32(video.ViewCount),
			LikeCount:     int32(video.LikeCount),
			FavoriteCount: int32(video.FavoriteCount),
			CommentCount:  int32(video.CommentCount),
		})
		if err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" { // Duplicate row video already exists
					continue
				}
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
		storeCount += 1
	}
	if storeCount > 0 {
		slog.Info("Stored Videos", "count", storeCount)
	}
	return nil
}

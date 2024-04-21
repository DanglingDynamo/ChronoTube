package repository

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

// Implements the VideoRepository interface
type YoutubeRepository struct {
	client *youtube.YoutubeClient
	db     *database.Queries
}

func NewYoutubeRepository(apiKey string, db *database.Queries) (*YoutubeRepository, error) {
	if apiKey != "" {
		return nil, errors.New("please provide an API key")
	}

	ytClient, err := youtube.NewYoutubeClient(context.Background(), os.Getenv("YOUTUBE_API_KEY"))
	if err != nil {
		return nil, err
	}

	return &YoutubeRepository{
		client: ytClient,
		db:     db,
	}, nil
}

func (repo *YoutubeRepository) FetchVideosFromAPI(
	query string,
	publishedAfter time.Time,
) ([]*models.Video, error) {
	slog.Info("Fetching Videos")
	ytVideos, err := repo.client.FetchVideos(
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

	videoStatistics, err := repo.client.FetchVideoStatistics(ids...)
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
			ViewCount:     videoStatistics[video.Id.VideoId].ViewCount,
			LikeCount:     videoStatistics[video.Id.VideoId].LikeCount,
			FavoriteCount: videoStatistics[video.Id.VideoId].FavoriteCount,
			CommentCount:  videoStatistics[video.Id.VideoId].CommentCount,
		}
	}

	return videos, nil
}

// Storing Videos is specific to the service hence implemented here
func (repo *YoutubeRepository) StoreVideos(ctx context.Context, videos []*models.Video) error {
	var pgErr *pgconn.PgError
	storeCount := 0
	for i := range videos {
		video := videos[i]
		_, err := repo.db.InsertVideo(ctx, database.InsertVideoParams{
			Title:         video.Title,
			Description:   video.Description,
			PublishedOn:   video.PublishedOn,
			ThumbnailUrl:  video.ThumbnailURL,
			Provider:      string(video.Provider),
			VideoID:       video.VideoID,
			ViewCount:     int64(video.ViewCount),
			LikeCount:     int64(video.LikeCount),
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
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

func (repo *YoutubeRepository) FetchPaginatedVideos(
	ctx context.Context,
	nextPage int,
) ([]models.Video, error) {
	offset := 5 * nextPage

	response, err := repo.db.FetchVideosPaginated(ctx, int32(offset))
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}

	videos := make([]models.Video, 0, 5) // Page size is 5
	for _, video := range response {
		videos = append(videos, models.VideoFromDatabaseVideo(video))
	}

	return videos, nil
}

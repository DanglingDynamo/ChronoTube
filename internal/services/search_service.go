package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/repository"
	"github.com/DanglingDynamo/chronotube/internal/utils"
)

type SearchService interface {
	QueryVideos(
		ctx context.Context,
		request models.PaginatedVideoRequest,
	) ([]models.Video, string, error)
}

type VideoSearchService struct {
	db repository.VideoRepository
}

func NewVideoSearchService(repo repository.VideoRepository) *VideoSearchService {
	return &VideoSearchService{
		db: repo,
	}
}

func (service *VideoSearchService) QueryVideos(
	ctx context.Context,
	request models.PaginatedVideoRequest,
) ([]models.Video, string, error) {
	nextPage := 0
	if request.NextPage != "" {
		decoded, err := utils.Decrypt(request.NextPage)
		if err != nil {
			return nil, "", err
		}

		nextPage, err = strconv.Atoi(string(decoded))
		if err != nil {
			return nil, "", err
		}
	}

	videos, err := service.db.FetchPaginatedVideos(ctx, nextPage)
	if err != nil {
		return nil, "", err
	}

	encoded := ""
	if len(videos) > 0 {
		var err error
		encoded, err = utils.Encrypt([]byte(fmt.Sprint(nextPage + 1)))
		if err != nil {
			return nil, "", err
		}
	}

	return videos, encoded, nil
}

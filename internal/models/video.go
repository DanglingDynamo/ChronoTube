package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/DanglingDynamo/chronotube/internal/database"
)

type Video struct {
	ID            uuid.UUID     `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	PublishedOn   time.Time     `json:"published_on"`
	ThumbnailURL  string        `json:"thumnail_url"`
	Provider      VideoProvider `json:"provider"`
	VideoID       string        `json:"video_key"`
	ViewCount     uint64        `json:"view_count"`
	LikeCount     uint64        `json:"like_count"`
	FavoriteCount uint64        `json:"favorite_count"`
	CommentCount  uint64        `json:"comment_count"`
}

func VideoFromDatabaseVideo(video database.Video) Video {
	return Video{
		ID:            video.ID,
		Title:         video.Title,
		Description:   video.Description,
		PublishedOn:   video.PublishedOn,
		ThumbnailURL:  video.ThumbnailUrl,
		Provider:      VideoProvider(video.Provider),
		VideoID:       video.VideoID,
		ViewCount:     uint64(video.ViewCount),
		LikeCount:     uint64(video.LikeCount),
		FavoriteCount: uint64(video.FavoriteCount),
		CommentCount:  uint64(video.CommentCount),
	}
}

type PaginatedVideoRequest struct {
	NextPage string `json:"next_page"`
}

type VideoProvider string

const (
	VideoProviderYoutube VideoProvider = "youtube"
)

package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID            uuid.UUID     `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	PublishedOn   time.Time     `json:"published_on"`
	ThumbnailURL  string        `json:"thumnail_url"`
	Provider      VideoProvider `json:"provider"`
	VideoID       string        `json:"video_key"`
	ViewCount     int           `json:"view_count"`
	LikeCount     int           `json:"like_count"`
	FavoriteCount int           `json:"favorite_count"`
	CommentCount  int           `json:"comment_count"`
}

type VideoProvider string

const (
	VideoProviderYoutube VideoProvider = "youtube"
)

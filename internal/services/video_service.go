package services

import "time"

type VideoService interface {
	FetchVideos(string, time.Time)
}

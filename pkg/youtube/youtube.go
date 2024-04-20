package youtube

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Creates a YoutubeClient that can be used to fetch videos
// It is made separately so that we can decouple the code from the third party
type YoutubeClient struct {
	service *youtube.Service
}

func NewYoutubeClient(ctx context.Context, apiKey string) (*YoutubeClient, error) {
	yt, err := youtube.NewService(
		ctx,
		option.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, err
	}

	return &YoutubeClient{
		service: yt,
	}, nil
}

func (client *YoutubeClient) FetchVideos(query string, publishedAfter time.Time) error {
	call := client.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		PublishedAfter(publishedAfter.Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		return err
	}

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			slog.Info(
				item.Snippet.Title,
				"url",
				"https://www.youtube.com/watch?v="+item.Id.VideoId,
				"thumbnail",
				item.Snippet.Thumbnails.Default.Url,
			)
		case "youtube#channel":
		case "youtube#playlyst":
			continue
		}
	}

	return nil
}

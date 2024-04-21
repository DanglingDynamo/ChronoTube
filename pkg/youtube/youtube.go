package youtube

import (
	"context"
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

func (client *YoutubeClient) FetchVideos(
	query string,
	publishedAfter time.Time,
) ([]*youtube.SearchResult, error) {
	call := client.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		PublishedAfter(publishedAfter.Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	videos := make(
		[]*youtube.SearchResult,
		0,
		len(response.Items),
	) // max videos can be equal to length
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos = append(videos, item)
		case "youtube#channel":
		case "youtube#playlyst":
			continue
		}
	}

	return videos, nil
}

func (client *YoutubeClient) FetchVideoStatistics(
	videoIDs ...string,
) (map[string]*youtube.VideoStatistics, error) {
	call := client.service.Videos.List([]string{"id", "snippet", "statistics"}).Id(videoIDs...)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	videoDetails := make(map[string]*youtube.VideoStatistics, len(response.Items))
	for _, item := range response.Items {
		videoDetails[item.Id] = item.Statistics
	}

	return videoDetails, nil
}

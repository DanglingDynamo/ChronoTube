package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	cronjobs "github.com/DanglingDynamo/chronotube/internal/cron_jobs"
	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/services"
)

func init() {
	godotenv.Load()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := services.NewYoutubeService(os.Getenv("API_SECRET_KEY"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	out := make(chan []*models.Video)

	go cronjobs.FetchVideos(ctx, time.Second*10, service, "cricket", out)

	go func() {
		for videos := range out {
			for i := range videos {
				slog.Info("video fetched", "title", videos[i].Title)
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	close(out)
}

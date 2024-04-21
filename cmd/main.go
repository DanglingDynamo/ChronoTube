package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/DanglingDynamo/chronotube/internal/config"
	cronjobs "github.com/DanglingDynamo/chronotube/internal/cron_jobs"
	"github.com/DanglingDynamo/chronotube/internal/database"
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
	// TODO: Cleanup
	config := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPass,
		config.DBName,
	)

	conn, err := sql.Open("pgx", uri)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	db := database.New(conn)

	service, err := services.NewYoutubeService(os.Getenv("API_SECRET_KEY"), db)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	out := make(chan []*models.Video)
	errChan := make(chan error)

	go cronjobs.FetchVideos(ctx, time.Second*10, service, "basketball", out, errChan)
	go cronjobs.StoreVideos(out, service, errChan)

	go func() {
		for err := range errChan {
			slog.Error(err.Error())
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	close(out)
	close(errChan)
}

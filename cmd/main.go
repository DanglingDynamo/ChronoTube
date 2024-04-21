package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/DanglingDynamo/chronotube/internal/config"
	"github.com/DanglingDynamo/chronotube/internal/constants"
	cronjobs "github.com/DanglingDynamo/chronotube/internal/cron_jobs"
	"github.com/DanglingDynamo/chronotube/internal/database"
	"github.com/DanglingDynamo/chronotube/internal/handlers"
	"github.com/DanglingDynamo/chronotube/internal/initializers"
	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/repository"
	"github.com/DanglingDynamo/chronotube/internal/routes"
	"github.com/DanglingDynamo/chronotube/internal/services"
)

func init() {
	godotenv.Load()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))
}

func main() {
	config := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan []*models.Video)
	errChan := make(chan error)
	initializers.InitDB(config.DBConfig)

	youtubeRepository, err := repository.NewYoutubeRepository(
		os.Getenv("API_SECRET_KEY"),
		database.New(initializers.DB),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	app := chi.NewRouter()

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: app,
	}

	router := routes.Router{
		SearchHandler: handlers.NewSearchHandler(
			services.NewVideoSearchService(youtubeRepository),
		),
	}
	routes.SetupRoutes(app, router)

	go cronjobs.FetchVideos(ctx, time.Second*10, youtubeRepository, "basketball", out, errChan)
	go cronjobs.StoreVideos(out, youtubeRepository, errChan)

	// Log errors in the goroutines also handle theme here if later required in case of emergency exit etc
	go func() {
		for err := range errChan {
			if errors.Is(err, constants.ErrAPIKeysUsed) {
				slog.Info("Out of API keys stopping FetchVideos")
				cancel()
				close(out)
				close(errChan)
			}
			slog.Error(err.Error())
		}
	}()

	interrupt := make(chan os.Signal, 1)
	go func() {
		<-interrupt
		_ = server.Shutdown(ctx)
		cancel()
		close(out)
		close(errChan)
		initializers.DB.Close()
	}()
	signal.Notify(interrupt, os.Interrupt)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error(err.Error())
	}
}

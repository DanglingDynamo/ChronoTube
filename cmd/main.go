package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/DanglingDynamo/chronotube/internal/services"
)

func init() {
	godotenv.Load()
}

func main() {
	service, err := services.NewYoutubeService(os.Getenv("API_SECRET_KEY"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	service.FetchVideos("ASMR", time.Now().AddDate(0, 0, -1))
}

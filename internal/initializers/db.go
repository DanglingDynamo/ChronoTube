package initializers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/DanglingDynamo/chronotube/internal/config"
)

var DB *sql.DB

func InitDB(config config.DBConfig) {
	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPass,
		config.DBName,
	)

	var err error
	DB, err = sql.Open("pgx", uri)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

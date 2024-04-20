package config

import "os"

type HTTPConfig struct {
	Port string
}

type DBConfig struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string
}

type Config struct {
	DBConfig
	HTTPConfig
}

func LoadConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			DBHost: os.Getenv("POSTGRES_HOST"),
			DBUser: os.Getenv("POSTGRES_USER"),
			DBPass: os.Getenv("POSTGRES_PASSWORD"),
			DBName: os.Getenv("POSTGRES_DB"),
			DBPort: os.Getenv("POSTGRES_PORT"),
		},
		HTTPConfig: HTTPConfig{
			Port: os.Getenv("PORT"),
		},
	}
}

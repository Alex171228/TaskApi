package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, errors.New("Can't load .env")
	}
	cnfg := Config{
		HTTPAddr:        os.Getenv("HTTP_ADDR"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ShutdownTimeout: 5 * time.Second,
	}
	shotdownvalue := os.Getenv("SHUTDOWN_TIMEOUT")
	if shotdownvalue != "" {
		timeout, err := time.ParseDuration(shotdownvalue)
		if err != nil {
			return Config{}, errors.New("SHUTDOWN_TIMEOUT must be a duration")
		}
		if timeout <= 0 {
			return Config{}, errors.New("SHUTDOWN_TIMEOUT must be more than zero")
		}
		cnfg.ShutdownTimeout = timeout
	}
	if cnfg.HTTPAddr == "" {
		cnfg.HTTPAddr = ":8082"
	}
	if cnfg.DatabaseURL == "" {
		cnfg.DatabaseURL = "postgres://todo:todo@localhost:5433/todo?sslmode=disable"
	}

	if strings.TrimSpace(cnfg.HTTPAddr) == "" {
		return Config{}, errors.New("HTTP_ADDR must not be empty")
	}
	if strings.TrimSpace(cnfg.DatabaseURL) == "" {
		return Config{}, errors.New("DATABASE_URL must not be empty")
	}

	return cnfg, nil
}

// добавить конфигурации шатдауна

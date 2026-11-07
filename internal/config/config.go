package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (Config, error) {
	cnfg := Config{
		HTTPAddr:    os.Getenv("HTTP_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
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

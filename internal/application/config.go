package application

import (
	"log/slog"
	"os"
)

type Config struct {
	Level    slog.Level
	HTTPAddr string
}

func NewConfig() (*Config, error) {
	var level slog.Level

	rawLogLevel := os.Getenv("LOG_LEVEL")
	if err := level.UnmarshalText([]byte(rawLogLevel)); err != nil {
		level = slog.LevelInfo
	}

	httpAddr := os.Getenv("HTTP_ADDR")

	return &Config{
		Level:    level,
		HTTPAddr: httpAddr,
	}, nil
}

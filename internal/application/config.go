package application

import (
	"log/slog"
	"os"
)

type Config struct {
	Level slog.Level
}

func NewConfig() (*Config, error) {
	var level slog.Level

	rawLogLevel := os.Getenv("LOG_LEVEL")
	if err := level.UnmarshalText([]byte(rawLogLevel)); err != nil {
		level = slog.LevelInfo
	}

	return &Config{
		Level: level,
	}, nil
}

package main

import (
	"fmt"
	"log/slog"
	"os"
	"reddit-clone/internal/logger"

	"github.com/joho/godotenv"
)

type App struct {
	log *slog.Logger
}

func NewApp(log *slog.Logger) *App {
	return &App{log: log}
}

func (a *App) Run() {
	a.log.Info("Service running... ")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file, using system default values")
	}

	var level slog.Level

	rawLogLevel := os.Getenv("LOG_LEVEL")
	if err := level.UnmarshalText([]byte(rawLogLevel)); err != nil {
		level = slog.LevelInfo
	}

	log := logger.New(level)

	app := NewApp(log)

	app.Run()
}

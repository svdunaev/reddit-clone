package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"reddit-clone/internal/application"
	server "reddit-clone/internal/http"
	"reddit-clone/internal/logger"
	"reddit-clone/internal/storage/inmem"

	"github.com/joho/godotenv"
	"k8s.io/utils/clock"
)

type app struct {
	log *slog.Logger
}

func NewApp(log *slog.Logger) *app {
	return &app{log: log}
}

func (a *app) Run(ctx context.Context) {
	a.log.Info("Service running... ")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file, using system default values")
	}

	cfg, err := application.NewConfig()
	if err != nil {
		log.Fatal("invalid config")
	}

	log := logger.New(cfg.Level)

	app := NewApp(log)

	app.Run(context.Background())

	repo := inmem.New(clock.RealClock{})

	srv := server.New(log, repo)

	srv.Start(cfg.HttpAddr)
}

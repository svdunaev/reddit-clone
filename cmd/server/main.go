package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"reddit-clone/internal/application"
	createPostCommand "reddit-clone/internal/application/command/create_post"
	"reddit-clone/internal/logger"
	"reddit-clone/internal/storage/inmem"
	server "reddit-clone/internal/transport/http"
	createPostHTTP "reddit-clone/internal/transport/http/create_post"

	"github.com/joho/godotenv"
	"k8s.io/utils/clock"
)

type App struct {
	log                      *slog.Logger
	repo                     *inmem.Store
	createPostHandlerHTTP    *createPostHTTP.Handler
	createPostCommandHandler *createPostCommand.Handler
}

func NewApp(log *slog.Logger) *App {
	repo := inmem.New(clock.RealClock{})

	createPostCommand := createPostCommand.NewHandler(repo)
	createPostHTTP := createPostHTTP.NewHandler(createPostCommand)

	return &App{log: log, repo: repo, createPostHandlerHTTP: createPostHTTP, createPostCommandHandler: createPostCommand}
}

func (a *App) Run(_ context.Context) {
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

	srv := server.New(log, app.createPostHandlerHTTP)
	app.Run(context.Background())

	srv.Start(cfg.HTTPAddr)
}

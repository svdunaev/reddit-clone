package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"reddit-clone/internal/application"
	createPostCommand "reddit-clone/internal/application/command/create_post"
	"reddit-clone/internal/logger"
	"reddit-clone/internal/repository/post"
	httptransport "reddit-clone/internal/transport/http"
	createPostHTTP "reddit-clone/internal/transport/http/create_post"
	"syscall"

	"github.com/joho/godotenv"
	"k8s.io/utils/clock"
)

type App struct {
	log                      *slog.Logger
	repo                     *post.Store
	createPostHandlerHTTP    *createPostHTTP.Handler
	createPostCommandHandler *createPostCommand.Handler
}

func NewApp(ctx context.Context) (*App, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file, using system default values")
	}

	cfg, err := application.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config %w", err)
	}

	logger := logger.New(cfg.Level)

	repo := post.NewInMem(clock.RealClock{})

	createPostCommand := createPostCommand.NewHandler(repo)
	createPostHTTP := createPostHTTP.NewHandler(createPostCommand)

	return &App{log: logger, repo: repo, createPostHandlerHTTP: createPostHTTP, createPostCommandHandler: createPostCommand}, nil
}

func (a *App) Close() error { return nil }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := NewApp(ctx)
	if err != nil {
		return err
	}
	defer app.Close()

	// хендлеры передаются из App в роутер
	router := httptransport.New(app.log, app.createPostHandlerHTTP)
	return router.Run(ctx, ":8080") // bind/start + graceful — внутри транспорта
}

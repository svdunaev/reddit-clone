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
	deletePostCommand "reddit-clone/internal/application/command/delete_post"
	updatePostCommand "reddit-clone/internal/application/command/update_post"
	getPostQuery "reddit-clone/internal/application/query/get_post"
	listPostsQuery "reddit-clone/internal/application/query/list_posts"
	"reddit-clone/internal/logger"
	"reddit-clone/internal/repository/post"
	httptransport "reddit-clone/internal/transport/http"
	createPostHTTP "reddit-clone/internal/transport/http/create_post"
	deletePostHttp "reddit-clone/internal/transport/http/delete_post"
	getPostHttp "reddit-clone/internal/transport/http/get_post"
	listPostsHttp "reddit-clone/internal/transport/http/list_posts"
	updatePostHttp "reddit-clone/internal/transport/http/update_post"
	"syscall"

	"github.com/joho/godotenv"
	"k8s.io/utils/clock"
)

type App struct {
	log                      *slog.Logger
	repo                     *post.Store
	createPostHandlerHTTP    *createPostHTTP.Handler
	createPostCommandHandler *createPostCommand.Handler
	getPostHandlerHttp       *getPostHttp.Handler
	getPostQueryHandler      *getPostQuery.Handler
	listPostsHandlerHttp     *listPostsHttp.Handler
	listPostsQueryHandler    *listPostsQuery.Handler
	updatePostHandlerHttp    *updatePostHttp.Handler
	updatePostCommandHandler *updatePostCommand.Handler
	deletePostHandlerHttp    *deletePostHttp.Handler
	deletePostCommandHandler *deletePostCommand.Handler
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

	getPostQuery := getPostQuery.NewHandler(repo)
	getPostHttp := getPostHttp.NewHandler(getPostQuery)

	listPostsQuery := listPostsQuery.NewHandler(repo)
	listPostsHttp := listPostsHttp.NewHandler(listPostsQuery)

	updatePostCommand := updatePostCommand.NewHandler(repo)
	updatePostHttp := updatePostHttp.NewHandler(updatePostCommand)

	deletePostCommand := deletePostCommand.NewHandler(repo)
	deletePostHttp := deletePostHttp.NewHandler(deletePostCommand)

	return &App{
		log:                      logger,
		repo:                     repo,
		createPostHandlerHTTP:    createPostHTTP,
		createPostCommandHandler: createPostCommand,
		getPostHandlerHttp:       getPostHttp,
		getPostQueryHandler:      getPostQuery,
		listPostsHandlerHttp:     listPostsHttp,
		listPostsQueryHandler:    listPostsQuery,
		updatePostHandlerHttp:    updatePostHttp,
		updatePostCommandHandler: updatePostCommand,
		deletePostHandlerHttp:    deletePostHttp,
		deletePostCommandHandler: deletePostCommand,
	}, nil
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
	router := httptransport.New(
		app.log,
		app.createPostHandlerHTTP,
		app.getPostHandlerHttp,
		app.listPostsHandlerHttp,
		app.updatePostHandlerHttp,
		app.deletePostHandlerHttp,
	)
	return router.Run(ctx, ":8080") // bind/start + graceful — внутри транспорта
}

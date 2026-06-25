package httptransport

import (
	"context"
	"log/slog"
	"net/http"
	"reddit-clone/internal/helpers/middlewares"
	createPostHTTP "reddit-clone/internal/transport/http/create_post"
	deletePostHttp "reddit-clone/internal/transport/http/delete_post"
	getPostHttp "reddit-clone/internal/transport/http/get_post"
	listPostsHttp "reddit-clone/internal/transport/http/list_posts"
	updatePostHttp "reddit-clone/internal/transport/http/update_post"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	logger *slog.Logger
}

func New(
	logger *slog.Logger,
	createHandler *createPostHTTP.Handler,
	getPostHandler *getPostHttp.Handler,
	listPostsHandler *listPostsHttp.Handler,
	updatePostHandler *updatePostHttp.Handler,
	deletePostHandler *deletePostHttp.Handler,
) *Server {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderMiddleware)

	s := &Server{
		router: r,
		logger: logger,
	}

	r.Use(s.LoggerMiddleware)
	s.routes(
		createHandler,
		getPostHandler,
		listPostsHandler,
		updatePostHandler,
		deletePostHandler,
	)

	return s
}

func (s *Server) routes(
	createHandler *createPostHTTP.Handler,
	getPostHandler *getPostHttp.Handler,
	listPostsHandler *listPostsHttp.Handler,
	updatePostHandler *updatePostHttp.Handler,
	deletePostHandler *deletePostHttp.Handler,
) {
	s.router.Route("/api/posts", func(r chi.Router) {
		r.Post("/", createHandler.HandleCreatePost)
		r.Get("/", listPostsHandler.HandleListPosts)
		r.Get("/{id}", getPostHandler.HandleGetPost)
		r.Put("/{id}", updatePostHandler.HandleUpdatePost)
		r.Delete("/{id}", deletePostHandler.HandleDeletePost)
	})
	s.router.Get("/health", s.health)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (r *Server) Run(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: r.router, ReadHeaderTimeout: 5 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	r.logger.Info("http listening", "addr", addr)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

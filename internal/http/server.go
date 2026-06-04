package server

import (
	"log/slog"
	"net/http"
	createHandler "reddit-clone/internal/handler/create"
	getByIdHandler "reddit-clone/internal/handler/get_by_id"
	"reddit-clone/internal/helpers/middlewares"
	"reddit-clone/internal/storage/inmem"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	logger *slog.Logger
	repo   *inmem.Store
}

func New(logger *slog.Logger, repo *inmem.Store) *Server {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderMiddleware)

	s := &Server{
		router: r,
		logger: logger,
		repo:   repo,
	}

	r.Use(s.loggerMiddleware)
	s.routes()

	return s
}

func (s *Server) routes() {
	s.router.Get("/health", s.health)
	s.router.Post("/api/posts", createHandler.NewHandler(s.repo).HandleCreatePost)
	s.router.Get("/api/posts/{id}", getByIdHandler.NewHandler(s.repo).HandleGetPost)
}

func (s *Server) Start(addr string) {
	http.ListenAndServe(addr, s.router)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.logger.Info("new request", "method:", r.Method, "path: ", r.URL.Path)
		next.ServeHTTP(w, r)
		s.logger.Info("request completed", "method:", r.Method, "path: ", r.URL.Path, "duration: ", time.Since(start))
	})
}

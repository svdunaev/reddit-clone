package http

import (
	"log"
	"log/slog"
	"net/http"
	"reddit-clone/internal/helpers/middlewares"
	createPostHTTP "reddit-clone/internal/transport/http/create_post"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	logger *slog.Logger
}

func New(logger *slog.Logger, createHandler *createPostHTTP.Handler) *Server {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.JSONHeaderMiddleware)

	s := &Server{
		router: r,
		logger: logger,
	}

	r.Use(s.LoggerMiddleware)
	s.routes(createHandler)

	return s
}

func (s *Server) Start(addr string) {
	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) routes(createHandler *createPostHTTP.Handler) {
	s.router.Route("/api/posts", func(r chi.Router) {
		r.Post("/", createHandler.HandleCreatePost)
	})
	s.router.Get("/health", s.health)
	// s.router.Get("/api/posts/{id}", getByIdHandler.NewHandler(s.repo).HandleGetPost)
	// s.router.Get("/api/posts", getListHandler.NewHandler(s.repo).HandleGetList)
	// s.router.Delete("/api/posts/{id}", deletePostHandler.NewHandler(s.repo).HandleDeletePost)
	// s.router.Put("/api/posts/{id}", updatePostHandler.NewHandler(s.repo).HandleUpdatePost)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

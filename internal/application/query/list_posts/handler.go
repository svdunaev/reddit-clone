package list_posts

import (
	"context"
	domain "reddit-clone/internal/domain/post"
)

type Handler struct {
	repo PostRepository
}

type ErrorDetails struct {
	Code string
}

type Result struct {
	Posts *[]domain.Post
	Error *ErrorDetails
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context) (Result, error) {
	res, err := h.repo.List(ctx)

	if err != nil {
		return Result{
			Posts: nil,
			Error: &ErrorDetails{
				Code: "internal_error",
			},
		}, err
	}

	return Result{
		Posts: &res,
		Error: nil,
	}, nil
}

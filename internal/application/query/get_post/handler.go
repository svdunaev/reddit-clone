package get_post

import (
	"context"
	"errors"
	domain "reddit-clone/internal/domain/post"

	"github.com/google/uuid"
)

type Handler struct {
	repo PostRepository
}

type Query struct {
	Id uuid.UUID
}

type ErrorDetails struct {
	Code string
}

type Result struct {
	Post  *domain.Post
	Error *ErrorDetails
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	post, err := h.repo.GetById(ctx, query.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Result{
				Post: nil,
				Error: &ErrorDetails{
					Code: NotFoundErrorCode,
				},
			}, err
		}

		return Result{
			Post:  nil,
			Error: nil,
		}, err
	}

	return Result{
		Post:  &post,
		Error: nil,
	}, nil
}

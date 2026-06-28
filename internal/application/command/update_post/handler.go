package update_post

import (
	"context"
	"errors"
	domain "reddit-clone/internal/domain/post"

	"github.com/google/uuid"
)

type Command struct {
	Id     uuid.UUID
	Title  string
	Body   string
	Author string
}

type ErrorDetails struct {
	Code    string
	Details []domain.ValidationError
}

type Result struct {
	Post   *domain.Post
	Errors *ErrorDetails
}

type Handler struct {
	repo PostRepository
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	post, err := h.repo.GetById(ctx, cmd.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Result{
				Post: nil,
				Errors: &ErrorDetails{
					Code:    NotFoundErrorCode,
					Details: nil,
				},
			}, err
		}
	}

	post.Author = cmd.Author
	post.Body = cmd.Body
	post.Title = cmd.Title

	if errs, err := post.Validate(); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return Result{
				Post: nil,
				Errors: &ErrorDetails{
					Code:    ValidationErrorCode,
					Details: errs,
				},
			}, err
		}

		return Result{
			Post:   nil,
			Errors: nil,
		}, err
	}

	updatedPost, err := h.repo.Update(ctx, cmd.Id, post)
	if err != nil {
		return Result{
			Post:   nil,
			Errors: nil,
		}, err
	}

	return Result{
		Post:   &updatedPost,
		Errors: nil,
	}, nil
}

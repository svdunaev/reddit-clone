package create_post

import (
	"context"
	"errors"
	domain "reddit-clone/internal/domain/post"
)

type Command struct {
	Title         string
	Body          string
	Author        string
	SubredditName string
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
	input := domain.NewPost(cmd.Author, cmd.Title, cmd.Body, cmd.SubredditName)

	if errs, err := input.Validate(); err != nil {
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

	post, err := h.repo.Create(ctx, *input)
	if err != nil {
		return Result{
			Post:   nil,
			Errors: nil,
		}, err
	}

	return Result{
		Post:   &post,
		Errors: nil,
	}, nil
}

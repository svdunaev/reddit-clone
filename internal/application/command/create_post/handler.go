package create_post

import (
	"context"
	"errors"
	"fmt"
	domain "reddit-clone/internal/domain/post"
	"strings"
)

type Command struct {
	Title         string
	Body          string
	Author        string
	SubredditName string
}

type Handler struct {
	repo PostRepository
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*domain.Post, error) {
	input := domain.NewPost(cmd.Author, cmd.Title, cmd.Body, cmd.SubredditName)

	if errs, err := input.Validate(); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			messages := make([]string, len(errs))
			for i, e := range errs {
				messages[i] = fmt.Sprintf("%s: %s", e.Field, e.Reason)
			}
			return &domain.Post{}, fmt.Errorf("%w: %v", domain.ErrValidation, strings.Join(messages, ": "))
		}

		return &domain.Post{}, err
	}

	post, err := h.repo.Create(ctx, *input)
	if err != nil {
		return &domain.Post{}, err
	}

	return &post, nil
}

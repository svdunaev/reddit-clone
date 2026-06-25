package delete_post

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Command struct {
	Id uuid.UUID
}

type Handler struct {
	repo PostRepository
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	if cmd.Id == uuid.Nil {
		return errors.New("id is required")
	}

	if err := h.repo.Delete(ctx, cmd.Id); err != nil {
		return err
	}

	return nil
}

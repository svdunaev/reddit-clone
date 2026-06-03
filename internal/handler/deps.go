package handler

import (
	"context"
	domain "reddit-clone/internal/domain/post"

	"github.com/google/uuid"
)

// PostRepository defines the contract for Post persistance operations
type PostRepository interface {
	//Create persists a new Post and returns the created Post or an error.
	Create(ctx context.Context, input domain.Post) (domain.Post, error)

	//GetById retrieves a single Post by its ID or returns an error if not found.
	GetById(ctx context.Context, id uuid.UUID) (domain.Post, error)

	//Delete removes the Post identified by id and returns an error if the operation fails.
	Delete(ctx context.Context, id uuid.UUID) error

	//List retrieves all Posts and returns them as a slice or an error.
	List(ctx context.Context) ([]domain.Post, error)

	//Update modifies an existing Post identified by id using the provided input,
	//and returns the updated Post or an error.
	Update(ctx context.Context, id uuid.UUID, input domain.Post) (domain.Post, error)
}

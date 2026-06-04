package get_by_id

import (
	"context"
	domain "reddit-clone/internal/domain/post"

	"github.com/google/uuid"
)

// PostRepository defines the contract for Post persistance operations
type PostRepository interface {
	//GetById retrieves a single Post by its ID or returns an error if not found.
	GetById(ctx context.Context, id uuid.UUID) (domain.Post, error)
}

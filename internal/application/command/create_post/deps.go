package create_post

import (
	"context"
	domain "reddit-clone/internal/domain/post"
)

// PostRepository defines the contract for Post persistance operations
type PostRepository interface {
	//Create persists a new Post and returns the created Post or an error.
	Create(ctx context.Context, input domain.Post) (domain.Post, error)
}

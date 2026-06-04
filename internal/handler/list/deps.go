package list

import (
	"context"
	domain "reddit-clone/internal/domain/post"
)

type PostRepository interface {
	//List retrieves all Posts and returns them as a slice or an error.
	List(ctx context.Context) ([]domain.Post, error)
}

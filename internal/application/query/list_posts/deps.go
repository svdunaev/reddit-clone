package list_posts

import (
	"context"
	domain "reddit-clone/internal/domain/post"
)

type PostRepository interface {
	List(ctx context.Context) ([]domain.Post, error)
}

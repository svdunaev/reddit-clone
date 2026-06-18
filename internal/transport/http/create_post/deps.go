package create_post

import (
	"context"
	"reddit-clone/internal/application/command/create_post"
	domain "reddit-clone/internal/domain/post"
)

type CreateCommand interface {
	Handle(ctx context.Context, cmd create_post.Command) (*domain.Post, error)
}

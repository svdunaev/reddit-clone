package create_post

import (
	"context"
	"reddit-clone/internal/application/command/create_post"
)

type CreatePostCommandHandler interface {
	Handle(ctx context.Context, cmd create_post.Command) (create_post.Result, error)
}

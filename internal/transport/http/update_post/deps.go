package update_post

import (
	"context"
	"reddit-clone/internal/application/command/update_post"
)

type UpdatePostCommandHandler interface {
	Handle(ctx context.Context, cmd update_post.Command) (update_post.Result, error)
}

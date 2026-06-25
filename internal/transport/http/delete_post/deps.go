package delete_post

import (
	"context"
	"reddit-clone/internal/application/command/delete_post"
)

type DeletePostCommandHandler interface {
	Handle(ctx context.Context, cmd delete_post.Command) error
}

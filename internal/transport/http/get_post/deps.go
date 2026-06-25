package get_post

import (
	"context"
	"reddit-clone/internal/application/query/get_post"
)

type GetPostQueryHandler interface {
	Handle(ctx context.Context, query get_post.Query) (get_post.Result, error)
}

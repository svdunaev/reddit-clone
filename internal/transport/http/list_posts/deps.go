package list_posts

import (
	"context"
	"reddit-clone/internal/application/query/list_posts"
)

type ListPostsQueryHandler interface {
	Handle(ctx context.Context) (list_posts.Result, error)
}

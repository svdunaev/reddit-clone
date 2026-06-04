package update

import (
	"context"
	domain "reddit-clone/internal/domain/post"

	"github.com/google/uuid"
)

type PostRepository interface {
	//Update modifies an existing Post identified by id using the provided input,
	//and returns the updated Post or an error.
	Update(ctx context.Context, id uuid.UUID, input domain.Post) (domain.Post, error)
}

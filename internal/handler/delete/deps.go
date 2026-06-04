package delete

import (
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	//Delete removes the Post identified by id and returns an error if the operation fails.
	Delete(ctx context.Context, id uuid.UUID) error
}

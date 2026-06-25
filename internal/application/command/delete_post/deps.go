package delete_post

import (
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

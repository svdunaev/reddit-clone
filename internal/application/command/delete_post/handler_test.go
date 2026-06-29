package delete_post

import (
	"context"
	"errors"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	id := uuid.New()

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	err := handler.repo.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func TestNilId(t *testing.T) {
	id := uuid.Nil
	expectedErr := errors.New("id is required")

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().Delete(gomock.Any(), id).Return(expectedErr)

	err := handler.repo.Delete(context.Background(), id)
	assert.Error(t, err)
}

func TestNotFound(t *testing.T) {
	id := uuid.New()

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().Delete(gomock.Any(), id).Return(domain.ErrNotFound)

	err := handler.repo.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

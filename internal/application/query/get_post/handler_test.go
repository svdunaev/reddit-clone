package get_post

import (
	"context"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	id := uuid.New()
	post := domain.Post{
		Id: id,
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().GetById(gomock.Any(), id).Return(post, nil)

	result, err := handler.repo.GetById(context.Background(), id)
	assert.NoError(t, err)

	assert.Equal(t, post, result)
}

func TestNotFound(t *testing.T) {
	id := uuid.New()
	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().GetById(gomock.Any(), id).Return(domain.Post{}, domain.ErrNotFound)

	result, err := handler.repo.GetById(context.Background(), id)

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Equal(t, domain.Post{}, result)
}

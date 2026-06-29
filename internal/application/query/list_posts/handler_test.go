package list_posts

import (
	"context"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	posts := []domain.Post{
		{
			Id: uuid.New(),
		},
		{
			Id: uuid.New(),
		},
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().List(gomock.Any()).Return(posts, nil)

	result, err := handler.repo.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, posts, result)
}

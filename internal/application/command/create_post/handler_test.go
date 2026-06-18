package create_post

import (
	"context"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	cmd := Command{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	expectedPost := domain.Post{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().Create(gomock.Any(), expectedPost).Return(expectedPost, nil)

	post, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, *post)
}

func TestValidationErrHandler(t *testing.T) {
	cmd := Command{
		Author:        "",
		Title:         "",
		Body:          "",
		SubredditName: "da",
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

package update_post

import (
	"context"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	post := domain.Post{
		Id:            uuid.New(),
		Author:        "Michael Jordan",
		Title:         "Chicago Bulls",
		Body:          "GOAT",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	expectedPost := domain.Post{
		Id:            post.Id,
		Author:        "Bill Russel",
		Title:         "Celtics",
		Body:          "Actual GOAT",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	cmd := Command{
		Id:            post.Id,
		Author:        "Bill Russel",
		Title:         "Celtics",
		Body:          "Actual GOAT",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().GetById(gomock.Any(), post.Id).Return(post, nil)
	repo.EXPECT().Update(gomock.Any(), post.Id, expectedPost).Return(expectedPost, nil)

	result, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, *result.Post)
}

func TestValidationErrHandler(t *testing.T) {
	post := domain.Post{
		Id:            uuid.New(),
		Author:        "Michael Jordan",
		Title:         "Chicago Bulls",
		Body:          "GOAT",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	cmd := Command{
		Id:     post.Id,
		Author: "",
		Title:  "",
		Body:   "",
	}

	expectedErrors := []domain.ValidationError{
		{
			Field:  "title",
			Reason: "title must be at least 3 characters but can not be more than 200",
		},
		{
			Field:  "author",
			Reason: "author can not be empty",
		},
		{
			Field:  "body",
			Reason: "body must be at least 1 character but can not be more than 10000",
		},
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)

	repo.EXPECT().GetById(gomock.Any(), post.Id).Return(post, nil)

	result, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, domain.ErrValidation, err)
	assert.NotNil(t, result.Errors)
	assert.Equal(t, ValidationErrorCode, result.Errors.Code)
	assert.Equal(t, expectedErrors, result.Errors.Details)
}

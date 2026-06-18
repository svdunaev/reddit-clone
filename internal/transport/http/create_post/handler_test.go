package create_post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	createPostCommand "reddit-clone/internal/application/command/create_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	post := &domain.Post{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	inputCommand := createPostCommand.Command{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	ctrl := gomock.NewController(t)
	cmd := NewMockCreateCommand(ctrl)
	handler := NewHandler(cmd)

	cmd.EXPECT().Handle(gomock.Any(), inputCommand).Return(post, nil)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(Request{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	})
	req, err := http.NewRequest(http.MethodPost, "/api/posts", &buf)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp domain.Post

	err = json.NewDecoder(rr.Body).Decode(&resp)

	assert.NoError(t, err)
	assert.Equal(t, *post, resp)
}

func TestBadJSONHandler(t *testing.T) {
	body := `{"author":"john","title":"test"`
	httpErr := utils.HttpError{
		Code: "bad_request",
	}

	ctrl := gomock.NewController(t)
	cmd := NewMockCreateCommand(ctrl)
	handler := NewHandler(cmd)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	handler.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp utils.HttpError
	json.NewDecoder(rr.Body).Decode(&resp)

	assert.Equal(t, httpErr.Code, resp.Code)
}

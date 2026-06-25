package create_post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	createPostCommand "reddit-clone/internal/application/command/create_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	expectedResult := createPostCommand.Result{
		Post: &domain.Post{
			Author:        "John Doe",
			Title:         "Aboba322",
			Body:          "ya ochen lublu testi",
			SubredditName: "https://pkg.go.dev/testing",
		},
		Errors: nil,
	}

	inputCommand := createPostCommand.Command{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing",
	}

	ctrl := gomock.NewController(t)
	cmd := NewMockCreatePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Post("/api/posts", handler.HandleCreatePost)

	cmd.EXPECT().Handle(gomock.Any(), inputCommand).Return(expectedResult, nil)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(Request{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing",
	})
	req, err := http.NewRequest(http.MethodPost, "/api/posts", &buf)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp respond.Post

	err = json.NewDecoder(rr.Body).Decode(&resp)

	assert.NoError(t, err)
	assert.Equal(t, respond.FromPost(expectedResult.Post), resp)
}

func TestBadJSONHandler(t *testing.T) {
	body := `{"author":"john","title":"test"`

	ctrl := gomock.NewController(t)
	cmd := NewMockCreatePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Post("/api/posts", handler.HandleCreatePost)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp respond.HttpError
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, "bad_request", resp.Error.Code)
}

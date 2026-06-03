package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/mocks"
	"reddit-clone/internal/storage/inmem"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappyCreateHandler(t *testing.T) {
	post := domain.Post{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	s := inmem.New(mockClock)
	now := time.Now()
	mockClock.EXPECT().Now().Return(now).Times(2)

	h := NewHandler(s)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	req, err := http.NewRequest(http.MethodPost, "/api/posts", &buf)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp domain.Post

	err = json.NewDecoder(rr.Body).Decode(&resp)

	assert.NoError(t, err)

	assert.Equal(t, post.Author, resp.Author)
	assert.Equal(t, post.Title, resp.Title)
	assert.Equal(t, post.Body, resp.Body)
	assert.Equal(t, post.SubredditName, resp.SubredditName)
}

func TestValidationErrCreateHandler(t *testing.T) {
	post := domain.Post{
		Author:        "1",
		Title:         "1",
		Body:          "",
		SubredditName: "123",
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	s := inmem.New(mockClock)

	h := NewHandler(s)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	h.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]any

	json.NewDecoder(rr.Body).Decode(&resp)

	errObj, ok := resp["error"].(map[string]any)
	assert.True(t, ok)

	assert.Equal(t, "validation_error", errObj["code"])
	assert.Equal(t, "validation failed", errObj["message"])

}

func TestBadJSONCreateHandler(t *testing.T) {
	body := `{"author":"john","title":"test"`

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	s := inmem.New(mockClock)

	h := NewHandler(s)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	h.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]any

	json.NewDecoder(rr.Body).Decode(&resp)

	errObj, ok := resp["error"].(map[string]any)
	assert.True(t, ok)

	assert.Equal(t, "bad_request", errObj["code"])

	msg, ok := errObj["message"].(string)
	assert.True(t, ok)
	assert.Contains(t, msg, "cannot unmarshal")
}

func TestNotFoundGetPostHandler(t *testing.T) {
	id := uuid.New().String()
	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	s := inmem.New(mockClock)
	h := NewHandler(s)

	r := chi.NewRouter()
	r.Get("/api/posts/{id}", h.HandleGetPost)

	req, _ := http.NewRequest(http.MethodGet, "/api/posts/"+id, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	var resp map[string]any

	json.NewDecoder(rr.Body).Decode(&resp)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	errObj, ok := resp["error"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "not_found", errObj["code"])

	msg, ok := errObj["message"].(string)
	assert.True(t, ok)
	assert.Contains(t, msg, "not found")
}

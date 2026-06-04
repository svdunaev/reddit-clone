package create

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	post := domain.Post{
		Author:        "John Doe",
		Title:         "Aboba322",
		Body:          "ya ochen lublu testi",
		SubredditName: "https://pkg.go.dev/testing#B.Helper",
	}

	ctrl := gomock.NewController(t)
	s := NewMockPostRepository(ctrl)
	h := NewHandler(s)

	s.EXPECT().
		Create(gomock.Any(), post).
		Return(post, nil)

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
	assert.Equal(t, post, resp)
}

func TestValidationErrHandler(t *testing.T) {
	post := domain.Post{
		Author:        "1",
		Title:         "1",
		Body:          "",
		SubredditName: "123",
	}

	httpErr := utils.HttpError{
		Code:    "validation_error",
		Message: "validation error",
	}

	ctrl := gomock.NewController(t)
	s := NewMockPostRepository(ctrl)
	h := NewHandler(s)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	h.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp utils.HttpError
	json.NewDecoder(rr.Body).Decode(&resp)

	assert.Equal(t, httpErr.Code, resp.Code)
	assert.Equal(t, httpErr.Message, resp.Message)
}

func TestBadJSONHandler(t *testing.T) {
	body := `{"author":"john","title":"test"`
	httpErr := utils.HttpError{
		Code: "bad_request",
	}

	ctrl := gomock.NewController(t)
	s := NewMockPostRepository(ctrl)
	h := NewHandler(s)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts", &buf)

	rr := httptest.NewRecorder()

	h.HandleCreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp utils.HttpError
	json.NewDecoder(rr.Body).Decode(&resp)

	assert.Equal(t, httpErr.Code, resp.Code)
}

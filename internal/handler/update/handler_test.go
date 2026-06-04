package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers/utils"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	id := uuid.New()
	post := domain.Post{
		Id:     id,
		Author: "Aboba",
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Put("/api/posts/{id}", handler.HandleUpdatePost)

	repo.EXPECT().Update(gomock.Any(), id, post).Return(post, nil)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)

	req, _ := http.NewRequest(http.MethodPut, "/api/posts/"+id.String(), &buf)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var res domain.Post
	err := json.NewDecoder(recorder.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, post, res)
}

func TestNotFound(t *testing.T) {
	id := uuid.Nil
	post := domain.Post{
		Id:     id,
		Author: "Aboba",
	}
	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	r := chi.NewRouter()
	r.Put("/api/posts/{id}", handler.HandleUpdatePost)

	repo.EXPECT().Update(gomock.Any(), id, post).Return(domain.Post{}, domain.ErrNotFound)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)

	req, _ := http.NewRequest(http.MethodPut, "/api/posts/"+id.String(), &buf)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestBadRequest(t *testing.T) {
	id := "aoaoao"

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	r := chi.NewRouter()
	r.Put("/api/posts/{id}", handler.HandleUpdatePost)

	req, _ := http.NewRequest(http.MethodPut, "/api/posts/"+id, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBadPayload(t *testing.T) {
	id := uuid.Nil
	post := "aoaoaoao"

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	r := chi.NewRouter()
	r.Put("/api/posts/{id}", handler.HandleUpdatePost)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)

	req, _ := http.NewRequest(http.MethodPut, "/api/posts/"+id.String(), &buf)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestValidationErrHandler(t *testing.T) {
	id := uuid.New()
	post := domain.Post{
		Author:        "1",
		Title:         "1",
		Body:          "",
		SubredditName: "123",
	}

	httpErr := utils.HttpError{
		Code:    "validation_error",
		Message: "validation failed",
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Put("/api/posts/{id}", handler.HandleUpdatePost)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(post)
	req, _ := http.NewRequest(http.MethodPut, "/api/posts/"+id.String(), &buf)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp utils.HttpError
	json.NewDecoder(recorder.Body).Decode(&resp)

	assert.Equal(t, httpErr.Code, resp.Code)
	assert.Equal(t, httpErr.Message, resp.Message)
}

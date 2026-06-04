package delete

import (
	"net/http"
	"net/http/httptest"
	domain "reddit-clone/internal/domain/post"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	id := uuid.New()

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/posts/"+id.String(), nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestNotFound(t *testing.T) {
	id := uuid.Nil

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	repo.EXPECT().Delete(gomock.Any(), id).Return(domain.ErrNotFound)

	req, _ := http.NewRequest(http.MethodDelete, "/api/posts/"+id.String(), nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestBadRequest(t *testing.T) {
	id := "eblan"

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	req, _ := http.NewRequest(http.MethodDelete, "/api/posts/"+id, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

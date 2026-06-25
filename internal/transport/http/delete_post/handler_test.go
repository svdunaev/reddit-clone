package delete_post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reddit-clone/internal/application/command/delete_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	id := uuid.New()

	ctrl := gomock.NewController(t)
	cmd := NewMockDeletePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	inputCommand := delete_post.Command{
		Id: id,
	}

	cmd.EXPECT().Handle(gomock.Any(), inputCommand).Return(nil)

	request, err := http.NewRequest(http.MethodDelete, "/api/posts/"+id.String(), nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestInvalidIdHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	cmd := NewMockDeletePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	invalidId := "sdflkm1231"

	req, err := http.NewRequest(http.MethodDelete, "/api/posts/"+invalidId, nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp respond.HttpError
	json.NewDecoder(recorder.Body).Decode(&resp)
	assert.Equal(t, IncorrectIdErrorCode, resp.Error.Code)
}

func TestNotFoundHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	cmd := NewMockDeletePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Delete("/api/posts/{id}", handler.HandleDeletePost)

	id := uuid.New()

	inputCmd := delete_post.Command{
		Id: id,
	}

	cmd.EXPECT().Handle(gomock.Any(), inputCmd).Return(domain.ErrNotFound)

	req, err := http.NewRequest(http.MethodDelete, "/api/posts/"+id.String(), nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var resp respond.HttpError
	json.NewDecoder(recorder.Body).Decode(&resp)
	assert.Equal(t, NotFoundErrorCode, resp.Error.Code)
}

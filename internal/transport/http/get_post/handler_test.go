package get_post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	getPostQuery "reddit-clone/internal/application/query/get_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHappyHandler(t *testing.T) {
	id := uuid.New()
	expectedResult := getPostQuery.Result{
		Post: &domain.Post{
			Id: id,
		},
		Error: nil,
	}

	inputQuery := getPostQuery.Query{
		Id: id,
	}

	ctrl := gomock.NewController(t)
	cmd := NewMockGetPostQueryHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Get("/api/posts/{id}", handler.HandleGetPost)

	cmd.EXPECT().Handle(gomock.Any(), inputQuery).Return(expectedResult, nil)

	req, err := http.NewRequest(http.MethodGet, "/api/posts/"+id.String(), nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp respond.Post
	err = json.NewDecoder(recorder.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respond.FromPost(expectedResult.Post), resp)
}

func TestInvalidIdHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	cmd := NewMockGetPostQueryHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Get("/api/posts/{id}", handler.HandleGetPost)

	invalidId := "sdflkm1231"

	req, err := http.NewRequest(http.MethodGet, "/api/posts/"+invalidId, nil)
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
	cmd := NewMockGetPostQueryHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Get("/api/posts/{id}", handler.HandleGetPost)

	id := uuid.New()

	inputQuery := getPostQuery.Query{
		Id: id,
	}

	expectedResult := getPostQuery.Result{
		Post: nil,
		Error: &getPostQuery.ErrorDetails{
			Code: NotFoundErrorCode,
		},
	}

	cmd.EXPECT().Handle(gomock.Any(), inputQuery).Return(expectedResult, domain.ErrNotFound)

	req, err := http.NewRequest(http.MethodGet, "/api/posts/"+id.String(), nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var resp respond.HttpError
	json.NewDecoder(recorder.Body).Decode(&resp)
	assert.Equal(t, NotFoundErrorCode, resp.Error.Code)
}

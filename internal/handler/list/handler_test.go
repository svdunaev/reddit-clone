package list

import (
	"encoding/json"
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
	posts := []domain.Post{
		{
			Id: uuid.New(),
		},
		{
			Id: uuid.New(),
		},
	}

	ctrl := gomock.NewController(t)
	repo := NewMockPostRepository(ctrl)
	handler := NewHandler(repo)
	router := chi.NewRouter()
	router.Get("/api/posts", handler.HandleGetList)

	repo.EXPECT().List(gomock.Any()).Return(posts, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/posts", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var res []domain.Post
	json.NewDecoder(recorder.Body).Decode(&res)
	assert.Equal(t, len(posts), len(res))
	assert.ElementsMatch(t, posts, res)
}

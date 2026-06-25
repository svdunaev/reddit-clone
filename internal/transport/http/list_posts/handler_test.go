package list_posts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	listPostsQuery "reddit-clone/internal/application/query/list_posts"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
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

	expectedResult := listPostsQuery.Result{
		Posts: &posts,
		Error: nil,
	}

	ctrl := gomock.NewController(t)
	query := NewMockListPostsQueryHandler(ctrl)
	handler := NewHandler(query)
	router := chi.NewRouter()
	router.Get("/api/posts/", handler.HandleListPosts)

	query.EXPECT().Handle(gomock.Any()).Return(expectedResult, nil)

	request, err := http.NewRequest(http.MethodGet, "/api/posts/", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp []respond.Post
	mappedPosts := make([]respond.Post, len(posts))
	for _, p := range posts {
		mappedPosts = append(mappedPosts, respond.FromPost(&p))
	}
	err = json.NewDecoder(recorder.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, mappedPosts, resp)
}

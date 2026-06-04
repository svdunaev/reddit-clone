package get_by_id

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
	id := uuid.New()
	post := domain.Post{
		Id: id,
	}

	ctrl := gomock.NewController(t)
	s := NewMockPostRepository(ctrl)
	h := NewHandler(s)
	r := chi.NewRouter()
	r.Get("/api/posts/{id}", h.HandleGetPost)

	s.EXPECT().GetById(gomock.Any(), id).Return(post, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/posts/"+id.String(), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res domain.Post
	json.NewDecoder(rr.Body).Decode(&res)
	assert.Equal(t, post, res)
}

func TestNotFound(t *testing.T) {
	id := uuid.New()
	ctrl := gomock.NewController(t)
	s := NewMockPostRepository(ctrl)
	h := NewHandler(s)

	r := chi.NewRouter()
	r.Get("/api/posts/{id}", h.HandleGetPost)

	s.EXPECT().GetById(gomock.Any(), id).Return(domain.Post{}, domain.ErrNotFound)

	req, _ := http.NewRequest(http.MethodGet, "/api/posts/"+id.String(), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

package update_post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reddit-clone/internal/application/command/update_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHappy(t *testing.T) {
	id := uuid.New()
	inputCommand := update_post.Command{
		Id:            id,
		Author:        "Magic Johnson",
		Title:         "Not a Goat",
		Body:          "Lalala lalala lalala",
		SubredditName: "https://pkg.go.dev/testing",
	}

	expectedOutput := update_post.Result{
		Post: &domain.Post{
			Id:            inputCommand.Id,
			Author:        "Magic Johnson",
			Title:         "Not a Goat",
			Body:          "Lalala lalala lalala",
			SubredditName: "https://pkg.go.dev/testing",
		},
		Errors: nil,
	}

	ctrl := gomock.NewController(t)
	cmd := NewMockUpdatePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Put("/api/posts/{id}", handler.HandleUpdatePost)

	cmd.EXPECT().Handle(gomock.Any(), inputCommand).Return(expectedOutput, nil)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(Request{
		Author:        "Magic Johnson",
		Title:         "Not a Goat",
		Body:          "Lalala lalala lalala",
		SubredditName: "https://pkg.go.dev/testing",
	})
	assert.NoError(t, err)

	request, err := http.NewRequest(http.MethodPut, "/api/posts/"+id.String(), &buf)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp respond.Post
	err = json.NewDecoder(recorder.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respond.FromPost(expectedOutput.Post), resp)
}

func TestInvalidIdHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	cmd := NewMockUpdatePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Get("/api/posts/{id}", handler.HandleUpdatePost)

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

func TestBadJSONHandler(t *testing.T) {
	id := uuid.New()
	body := `{"author":"john","title":"test"`

	ctrl := gomock.NewController(t)
	cmd := NewMockUpdatePostCommandHandler(ctrl)
	handler := NewHandler(cmd)
	router := chi.NewRouter()
	router.Post("/api/posts/{id}", handler.HandleUpdatePost)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/posts/"+id.String(), &buf)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp respond.HttpError
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, "bad_request", resp.Error.Code)
}

package respond

import (
	"encoding/json"
	"net/http"
	createPostCommand "reddit-clone/internal/application/command/create_post"
	domain "reddit-clone/internal/domain/post"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID `json:"id,omitempty"`
	Author        string    `json:"author"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	SubredditName string    `json:"subreddit_name"`
	Score         int64     `json:"score,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type Errors struct {
	Code    string                   `json:"code,omitempty"`
	Details []domain.ValidationError `json:"details,omitempty"`
}

type Response struct {
	Post   Post   `json:"post"`
	Errors Errors `json:"errors"`
}

func ToResponse(result createPostCommand.Result) Response {
	var resp Response

	if result.Post != nil {
		resp.Post = Post{
			Id:            result.Post.Id,
			Author:        result.Post.Author,
			Body:          result.Post.Body,
			Title:         result.Post.Title,
			SubredditName: result.Post.SubredditName,
			Score:         result.Post.Score,
			CreatedAt:     result.Post.CreatedAt,
			UpdatedAt:     result.Post.UpdatedAt,
		}
	}

	if result.Errors != nil {
		resp.Errors = Errors{
			Code:    result.Errors.Code,
			Details: result.Errors.Details,
		}
	}

	return resp
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, code, msg string) {
	JSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg}})
}

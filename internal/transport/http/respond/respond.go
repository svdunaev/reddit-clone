package respond

import (
	"encoding/json"
	"net/http"
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

type HttpError struct {
	Error Errors `json:"error"`
}

func FromPost(p *domain.Post) Post {
	return Post{
		Id: p.Id, Author: p.Author, Title: p.Title, Body: p.Body,
		SubredditName: p.SubredditName, Score: p.Score,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, code, msg string) {
	JSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg}})
}

func ValidationFailed(w http.ResponseWriter, code string, details []domain.ValidationError) {
	JSON(w, http.StatusBadRequest, map[string]any{"errors": map[string]any{"code": code, "details": details}})
}

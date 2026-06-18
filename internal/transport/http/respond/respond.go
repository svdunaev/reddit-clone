package respond

import (
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

func FromPost(p *domain.Post) Post {
	return Post{
		Id:            p.Id,
		Title:         p.Title,
		Body:          p.Body,
		Author:        p.Author,
		SubredditName: p.SubredditName,
		Score:         p.Score,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

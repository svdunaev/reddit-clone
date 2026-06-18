package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID
	Author        string
	Title         string
	Body          string
	SubredditName string
	Score         int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPost(author, title, body, subredditName string) *Post {
	return &Post{
		Author:        author,
		Title:         title,
		Body:          body,
		SubredditName: subredditName,
	}
}

func (p *Post) Validate() ([]ValidationError, error) {
	errors := make([]ValidationError, 0)

	if len(p.Title) < 3 || len(p.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:  "title",
			Reason: "title must be at least 3 characters but can not be more than 200",
		})
	}

	if len(p.Author) < 1 {
		errors = append(errors, ValidationError{
			Field:  "author",
			Reason: "author can not be empty",
		})
	}

	if len(p.Body) < 1 || len(p.Body) > 10000 {
		errors = append(errors, ValidationError{
			Field:  "body",
			Reason: "body must be at least 1 character but can not be more than 10000",
		})
	}

	_, err := url.ParseRequestURI(p.SubredditName)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:  "subredditName",
			Reason: "subreddit name is not a valid URL",
		})
	}

	if len(errors) != 0 {
		return errors, ErrValidation
	}

	return nil, nil
}

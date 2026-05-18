package domain

import (
	"errors"
	"net/url"
	"time"
)

type Post struct {
	Id            string
	Author        string
	Title         string
	SubredditName string
	Body          string
	Score         int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ValidationError struct {
	Field  string
	Reason string
}

var ErrValidation = errors.New("validation failed")

func (p *Post) Validate() ([]ValidationError, error) {
	errors := make([]ValidationError, 0)

	if len(p.Title) <= 3 || len(p.Title) > 200 {
		errors = append(errors, ValidationError{
			Field:  p.Title,
			Reason: "title must be at least 3 characters but can not be more than 200",
		})
	}

	if len(p.Author) < 1 {
		errors = append(errors, ValidationError{
			Field:  p.Author,
			Reason: "author can not be empty",
		})
	}

	if len(p.Body) < 1 || len(p.Body) > 10000 {
		errors = append(errors, ValidationError{
			Field:  p.Body,
			Reason: "body must be at least 1 character but can not be more than 10000",
		})
	}

	_, err := url.ParseRequestURI(p.SubredditName)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:  p.SubredditName,
			Reason: "subreddit name is not a valid URL",
		})
	}

	return errors, nil
}

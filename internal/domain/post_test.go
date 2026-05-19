package domain

import "testing"

func TestPostValidate(t *testing.T) {
	tests := []struct {
		name          string
		post          Post
		wantErrCount  int
		wantErrFields []string
	}{
		{
			name: "happy path",
			post: Post{
				Title:         "Valid Title",
				Author:        "author123",
				Body:          "This is a valid body.",
				SubredditName: "https://go.dev/ref/spec",
			},
			wantErrCount:  0,
			wantErrFields: []string{},
		},
		{
			name: "invalid title: too short",
			post: Post{
				Title:         "hi",
				Author:        "author123",
				Body:          "This is a valid body.",
				SubredditName: "https://go.dev/ref/spec",
			},
			wantErrCount:  1,
			wantErrFields: []string{"title"},
		},
		{
			name: "invalid author: empty author",
			post: Post{
				Title:         "Valid title",
				Author:        "",
				Body:          "This is a valid body.",
				SubredditName: "https://go.dev/ref/spec",
			},
			wantErrCount:  1,
			wantErrFields: []string{"author"},
		},
		{
			name: "invalid body: empty body",
			post: Post{
				Title:         "Valid title",
				Author:        "author123",
				Body:          "",
				SubredditName: "https://go.dev/ref/spec",
			},
			wantErrCount:  1,
			wantErrFields: []string{"body"},
		},
		{
			name: "invalid subreddit name: not a valid url",
			post: Post{
				Title:         "Valid title",
				Author:        "author123",
				Body:          "author123",
				SubredditName: "not-a-url",
			},
			wantErrCount:  1,
			wantErrFields: []string{"subredditName"},
		},
		{
			name: "all fields are invalid",
			post: Post{
				Title:         "hi",
				Author:        "",
				Body:          "",
				SubredditName: "not-a-url",
			},
			wantErrCount:  4,
			wantErrFields: []string{"title", "author", "body", "subredditName"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs, err := tt.post.Validate()
			if err != nil {
				t.Fatalf("unexpected error returned %v\n", err)
			}

			if len(errs) != tt.wantErrCount {
				t.Errorf("expected %d validation erros, got %d: %v\n", tt.wantErrCount, len(errs), errs)
			}

			for i, wantField := range tt.wantErrFields {
				if i >= len(errs) {
					t.Errorf("missing expected error for field: %q\n", wantField)
				}
				if errs[i].Field != wantField {
					t.Errorf("error[%d]: expected Field %q, got %q\n", i, wantField, errs[i].Field)
				}
			}
		})
	}
}

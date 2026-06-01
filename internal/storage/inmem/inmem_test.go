package inmem

import (
	"context"
	"fmt"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/mocks"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var fixedUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.Post
		wantErr bool
	}{
		{
			name: "successfully creates a post",
			input: domain.Post{
				Author: "test-author",
				Title:  "test-title",
				Body:   "test-body-test-body-test-body",
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClock := mocks.NewMockClock(ctrl)
	now := time.Now()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClock.EXPECT().Now().Return(now).Times(2)

			s := New(mockClock)
			got, err := s.Create(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.NotEqual(t, uuid.Nil, got.Id)
			assert.Equal(t, tt.input.Author, got.Author)
			assert.Equal(t, tt.input.Body, got.Body)
			assert.Equal(t, tt.input.Title, got.Title)
			assert.Equal(t, now, got.CreatedAt)
			assert.Equal(t, now, got.UpdatedAt)
		})

	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name      string
		seedPosts []domain.Post
		wantCount int
	}{
		{
			name:      "returns empty list when no posts",
			seedPosts: []domain.Post{},
			wantCount: 0,
		},
		{
			name: "returns all posts",
			seedPosts: []domain.Post{
				{Title: "Post 1", Author: "author1", Body: "Body 1"},
				{Title: "Post 2", Author: "author2", Body: "Body 2"},
			},
			wantCount: 2,
		},
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	now := time.Now()
	ctx := context.Background()
	s := New(mockClock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClock.EXPECT().Now().Return(now).Times(len(tt.seedPosts) * 2)

			for _, p := range tt.seedPosts {
				_, err := s.Create(ctx, p)
				assert.NoError(t, err)
			}

			got, err := s.List(ctx)

			assert.NoError(t, err)
			assert.Len(t, got, tt.wantCount)
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Post
		targetId func(targetId uuid.UUID) uuid.UUID
		wantErr  bool
	}{
		{
			name:     "successfully retrieves a post",
			input:    domain.Post{Title: "Test Title", Author: "author1", Body: "Test body"},
			targetId: func(targetId uuid.UUID) uuid.UUID { return targetId },
		},
		{
			name:     "returns error for non existent id",
			targetId: func(_ uuid.UUID) uuid.UUID { return fixedUUID },
			wantErr:  true,
		},
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	now := time.Now()
	s := New(mockClock)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				mockClock.EXPECT().Now().Return(now).Times(2)
			}

			var createdId uuid.UUID
			if !tt.wantErr {
				created, err := s.Create(ctx, tt.input)
				assert.NoError(t, err)
				createdId = created.Id
			}

			createdPost := domain.Post{
				Id:        createdId,
				Author:    tt.input.Author,
				Title:     tt.input.Title,
				Body:      tt.input.Body,
				CreatedAt: now,
				UpdatedAt: now,
			}

			post, err := s.GetById(ctx, createdId)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, createdPost, post)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Post
		deleteId func(createdId uuid.UUID) uuid.UUID
		wantErr  bool
	}{
		{
			name:     "successfully deletes existing post",
			input:    domain.Post{Title: "Test Title", Author: "author1", Body: "Test body"},
			deleteId: func(createdId uuid.UUID) uuid.UUID { return createdId },
		},
		{
			name:     "returns error for non existent id",
			deleteId: func(_ uuid.UUID) uuid.UUID { return fixedUUID },
			wantErr:  true,
		},
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	now := time.Now()
	s := New(mockClock)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				mockClock.EXPECT().Now().Return(now).Times(2)
			}

			var createdId uuid.UUID
			if !tt.wantErr {
				created, err := s.Create(ctx, tt.input)
				assert.NoError(t, err)
				createdId = created.Id
			}

			err := s.Delete(ctx, tt.deleteId(createdId))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			_, getErr := s.GetById(ctx, createdId)
			assert.Error(t, getErr)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		initialPost   domain.Post
		updateId      func(updateId uuid.UUID) uuid.UUID
		updatePayload domain.Post
		wantErr       bool
	}{
		{
			name: "successfully updates an entire post",
			initialPost: domain.Post{
				Author: "test-author",
				Title:  "test-title",
				Body:   "test-body-test-body-test-body",
			},
			updateId: func(updateId uuid.UUID) uuid.UUID { return updateId },
			updatePayload: domain.Post{
				Author: "updated-author",
				Title:  "updated-title",
				Body:   "updated-body-updated-body",
			},
			wantErr: false,
		},
		{
			name: "successfully updates a title of a post",
			initialPost: domain.Post{
				Author: "test-author",
				Title:  "test-title",
				Body:   "test-body-test-body-test-body",
			},
			updateId: func(updateId uuid.UUID) uuid.UUID { return updateId },
			updatePayload: domain.Post{
				Title: "updated-title",
			},
			wantErr: false,
		},
		{
			name: "successfully updates an author of a post",
			initialPost: domain.Post{
				Author: "test-author",
				Title:  "test-title",
				Body:   "test-body-test-body-test-body",
			},
			updateId: func(updateId uuid.UUID) uuid.UUID { return updateId },
			updatePayload: domain.Post{
				Author: "updated-author",
			},
			wantErr: false,
		},
		{
			name: "successfully updates a body of a post",
			initialPost: domain.Post{
				Author: "test-author",
				Title:  "test-title",
				Body:   "test-body-test-body-test-body",
			},
			updateId: func(updateId uuid.UUID) uuid.UUID { return updateId },
			updatePayload: domain.Post{
				Body: "updated-body-updated-body",
			},
			wantErr: false,
		},
		{
			name:     "returns error for non existent id",
			updateId: func(_ uuid.UUID) uuid.UUID { return fixedUUID },
			wantErr:  true,
		},
	}

	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	now := time.Now()
	s := New(mockClock)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				mockClock.EXPECT().Now().Return(now).Times(3)
			}

			var createdId uuid.UUID
			if !tt.wantErr {
				created, err := s.Create(ctx, tt.initialPost)
				assert.NoError(t, err)
				createdId = created.Id
			}

			post, err := s.Update(ctx, tt.updateId(createdId), tt.updatePayload)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, createdId, post.Id)
			assert.Equal(t, tt.updatePayload.Title, post.Title)
			assert.Equal(t, tt.updatePayload.Author, post.Author)
			assert.Equal(t, tt.updatePayload.Body, post.Body)
			assert.Equal(t, now, post.UpdatedAt)
		})
	}
}

func TestStore_Concurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Now()).AnyTimes()
	s := New(mockClock)
	goroutines := 10

	var wg sync.WaitGroup

	// Create
	createdPosts := make([]domain.Post, goroutines)
	createdIds := make([]uuid.UUID, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			post, err := s.Create(context.Background(), domain.Post{
				Author: fmt.Sprintf("test-author-%d", i),
				Title:  fmt.Sprintf("test-title-%d", i),
				Body:   fmt.Sprintf("test-body-%d", i),
			})
			assert.NoError(t, err)
			createdPosts[i] = post
		}(i)
	}

	wg.Wait()
	for i, p := range createdPosts {
		createdIds[i] = p.Id
	}

	// List

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			posts, err := s.List(context.Background())
			assert.NoError(t, err)
			assert.ElementsMatch(t, createdPosts, posts)
		}()
	}
	wg.Wait()

	// GetById

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			post, err := s.GetById(context.Background(), createdIds[i])
			assert.NoError(t, err)
			assert.Equal(t, createdPosts[i], post)
		}(i)
	}

	wg.Wait()

	// Update

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			updated, err := s.Update(context.Background(), createdIds[i], domain.Post{
				Author: fmt.Sprintf("updated-author-%d", i),
				Title:  fmt.Sprintf("updated-title-%d", i),
				Body:   fmt.Sprintf("updated-body-%d", i),
			})
			assert.NoError(t, err)
			assert.Equal(t, createdIds[i], updated.Id)
			assert.Equal(t, fmt.Sprintf("updated-title-%d", i), updated.Title)
			assert.Equal(t, fmt.Sprintf("updated-author-%d", i), updated.Author)
			assert.Equal(t, fmt.Sprintf("updated-body-%d", i), updated.Body)
		}(i)
	}
	wg.Wait()

	// Delete

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			err := s.Delete(context.Background(), createdIds[i])
			assert.NoError(t, err)

			_, err = s.GetById(context.Background(), createdIds[i])
			assert.Error(t, err)

		}(i)
	}
	wg.Wait()
}

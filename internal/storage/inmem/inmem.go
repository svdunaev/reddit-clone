package inmem

import (
	"context"
	"fmt"
	"maps"
	domain "reddit-clone/internal/domain/post"
	"sync"

	"github.com/google/uuid"
	"k8s.io/utils/clock"
)

type Store struct {
	mu     sync.RWMutex
	posts  map[uuid.UUID]*domain.Post
	nextId int64
	clock  clock.Clock
}

func New(clock clock.Clock) *Store {
	return &Store{
		posts: make(map[uuid.UUID]*domain.Post),
		clock: clock,
	}
}

func (s *Store) Create(ctx context.Context, input domain.Post) (domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return domain.Post{}, err
	}

	input.Id = uuid.New()
	input.CreatedAt = s.clock.Now()
	input.UpdatedAt = s.clock.Now()

	s.posts[input.Id] = &input

	return *s.posts[input.Id], nil
}

func (s *Store) List(ctx context.Context) ([]domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return []domain.Post{}, err
	}

	res := make([]domain.Post, 0, len(s.posts))

	for post := range maps.Values(s.posts) {
		res = append(res, *post)
	}

	return res, nil
}

func (s *Store) GetById(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return domain.Post{}, err
	}

	post, ok := s.posts[id]
	if !ok {
		return domain.Post{}, domain.ErrNotFound
	}

	return *post, nil
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	_, ok := s.posts[id]
	if !ok {
		return fmt.Errorf("post with id %d does not exist", id)
	}

	delete(s.posts, id)

	return nil
}

func (s *Store) Update(ctx context.Context, id uuid.UUID, input domain.Post) (domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return domain.Post{}, err
	}

	_, ok := s.posts[id]
	if !ok {
		return domain.Post{}, fmt.Errorf("post with id %d does not exist", id)
	}

	input.Id = id
	s.posts[id] = &input
	s.posts[id].UpdatedAt = s.clock.Now()

	return *s.posts[id], nil
}

package memory

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/example/iching-app/internal/domain"
)

type ReadingRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Reading
}

func NewReadingRepository() *ReadingRepository {
	return &ReadingRepository{items: map[string]domain.Reading{}}
}

func (r *ReadingRepository) Save(_ context.Context, reading domain.Reading) (domain.Reading, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[reading.ID] = reading
	return reading, nil
}

func (r *ReadingRepository) List(_ context.Context) ([]domain.Reading, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Reading, 0, len(r.items))
	for _, v := range r.items {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *ReadingRepository) Get(_ context.Context, id string) (domain.Reading, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.items[id]
	if !ok {
		return domain.Reading{}, errors.New("reading not found")
	}
	return v, nil
}

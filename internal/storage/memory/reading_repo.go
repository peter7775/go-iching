package memory

import (
	"context"
	"errors"
	"time"

	"github.com/example/iching-fiber-app/internal/domain"
)

type ReadingRepository struct {
	items []domain.Reading
}

func NewReadingRepository() *ReadingRepository {
	return &ReadingRepository{items: []domain.Reading{}}
}

func (r *ReadingRepository) Save(_ context.Context, in domain.Reading) (domain.Reading, error) {
	r.items = append([]domain.Reading{in}, r.items...)
	return in, nil
}

func (r *ReadingRepository) List(_ context.Context) ([]domain.Reading, error) {
	out := make([]domain.Reading, len(r.items))
	copy(out, r.items)
	return out, nil
}

func (r *ReadingRepository) Get(_ context.Context, id string) (domain.Reading, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Reading{}, errors.New("reading not found")
}

func (r *ReadingRepository) SaveReflection(_ context.Context, id string, rating int, note string, createdAt time.Time) error {
	for i := range r.items {
		if r.items[i].ID == id {
			r.items[i].Reflection.Rating = rating
			r.items[i].Reflection.Note = note
			r.items[i].Reflection.CreatedAt = &createdAt
			return nil
		}
	}
	return errors.New("reading not found")
}

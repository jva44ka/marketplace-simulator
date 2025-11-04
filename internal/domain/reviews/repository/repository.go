package repository

import (
	"context"
	"gitlab.ozon.dev/16/students/week-1-workshop/internal/domain/model"
	"sync"
	"sync/atomic"
)

type InMemoryRepository struct {
	storage map[model.Sku][]model.Review
	mx      sync.RWMutex

	idFactory atomic.Uint64
}

func NewReviewRepository(cap int) *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(map[model.Sku][]model.Review, cap),
	}
}

func (r *InMemoryRepository) CreateReview(_ context.Context, review model.Review) (model.Review, error) {
	reviewID := r.idFactory.Add(1)
	review.ID = model.ReviewID(reviewID)

	r.mx.Lock()
	defer r.mx.Unlock()
	r.storage[review.Sku] = append(r.storage[review.Sku], review)

	return review, nil
}

func (r *InMemoryRepository) GetReviewsBySku(_ context.Context, sku model.Sku) ([]model.Review, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.storage[sku], nil
}

package repository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/jva44ka/ozon-simulator-go/internal/domain/model"
)

type InMemoryProductRepository struct {
	storage map[uint64]model.Product
	mx      sync.RWMutex

	idFactory atomic.Uint64
}

func NewProductRepository(cap int) *InMemoryProductRepository {
	return &InMemoryProductRepository{
		storage: make(map[uint64]model.Product, cap),
	}
}

func (r *InMemoryProductRepository) GetProductBySku(_ context.Context, sku uint64) (*model.Product, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	product, productExists := r.storage[sku]

	if productExists == false {
		return nil, errors.New("Product not found in storage")
	}

	return &product, nil
}

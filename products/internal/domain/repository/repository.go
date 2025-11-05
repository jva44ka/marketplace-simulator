package repository

import (
	"context"
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
	storage := make(map[uint64]model.Product, cap)

	storage[1] = model.Product{
		Sku:   1,
		Price: 100.0,
		Name:  "Крем для лица"}
	storage[2] = model.Product{
		Sku:   2,
		Price: 600.0,
		Name:  "Дворники для лады весты"}
	storage[3] = model.Product{
		Sku:   3,
		Price: 600.0,
		Name:  "Вареники из Ozon Fresh"}

	return &InMemoryProductRepository{
		storage: storage,
	}
}

func (r *InMemoryProductRepository) GetProductBySku(_ context.Context, sku uint64) (*model.Product, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	product, productExists := r.storage[sku]

	if productExists == false {
		return nil, nil
	}

	return &product, nil
}

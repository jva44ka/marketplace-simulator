package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jva44ka/ozon-simulator-go/internal/domain/model"
)

type ProductRepository interface {
	GetProductBySku(_ context.Context, sku uint64) (*model.Product, error)
}

type ProductService struct {
	productRepository ProductRepository
}

func (s *ProductService) GetProductBySku(ctx context.Context, sku uint64) (*model.Product, error) {
	if sku < 1 {
		return nil, errors.New("sku must be passed")
	}

	product, err := s.productRepository.GetProductBySku(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("productRepository.GetProductBySku :%w", err)
	}

	return product, nil
}

func NewProductService(productRepository ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

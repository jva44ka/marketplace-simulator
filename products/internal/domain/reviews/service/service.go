package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/16/students/week-1-workshop/internal/domain/model"
)

type ReviewRepository interface {
	CreateReview(_ context.Context, review model.Review) (model.Review, error)
	GetReviewsBySku(_ context.Context, sku model.Sku) ([]model.Review, error)
}

type ProductService interface {
	GetProductBySku(ctx context.Context, sku model.Sku) (*model.Product, error)
}

type ReviewService struct {
	reviewRepository ReviewRepository
	productService   ProductService
}

func NewReviewService(reviewRepository ReviewRepository, productService ProductService) *ReviewService {
	return &ReviewService{reviewRepository: reviewRepository, productService: productService}
}

func (s *ReviewService) AddReview(ctx context.Context, review model.Review) (model.Review, error) {
	if review.Sku < 1 || review.UserID == uuid.Nil {
		return model.Review{}, errors.New("sku and user_id must be passed")
	}

	_, err := s.productService.GetProductBySku(ctx, review.Sku)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			return model.Review{}, fmt.Errorf("productService.GetProductBySku: %w", err)
		}

		return model.Review{}, err
	}

	newReview, err := s.reviewRepository.CreateReview(ctx, review)
	if err != nil {
		return model.Review{}, fmt.Errorf("reviewRepository.CreateReview :%w", err)
	}

	return newReview, nil
}

func (s *ReviewService) GetReviewsBySku(ctx context.Context, sku model.Sku) ([]model.Review, error) {
	if sku < 1 {
		return nil, errors.New("sku must be passed")
	}

	reviews, err := s.reviewRepository.GetReviewsBySku(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("reviewRepository.GetReviewsBySku :%w", err)
	}

	return reviews, nil
}

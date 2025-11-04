package get_reviews_by_sku_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/16/students/week-1-workshop/internal/domain/model"
	http2 "gitlab.ozon.dev/16/students/week-1-workshop/pkg/http"
	"net/http"
	"strconv"
)

type ReviewService interface {
	GetReviewsBySku(ctx context.Context, sku model.Sku) ([]model.Review, error)
}

type GetReviewsBySkuHandler struct {
	reviewService ReviewService
}

func NewGetReviewsBySkuHandler(reviewService ReviewService) *GetReviewsBySkuHandler {
	return &GetReviewsBySkuHandler{reviewService: reviewService}
}

func (h GetReviewsBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	skuRaw := r.PathValue("sku")
	sku, err := strconv.Atoi(skuRaw)
	if err != nil {
		if err = http2.ErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	if sku < 1 {
		if err = http2.ErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	reviews, err := h.reviewService.GetReviewsBySku(r.Context(), model.Sku(sku))
	if err != nil {
		if err = http2.ErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {
			return
		}

		return
	}

	response := GetReviewsResponse{Reviews: make([]GetReviewsReviewResponse, 0, len(reviews))}
	for _, review := range reviews {
		response.Reviews = append(response.Reviews, GetReviewsReviewResponse{
			ID:      uint64(review.ID),
			Sku:     uint64(review.Sku),
			Comment: review.Comment,
			UserID:  review.UserID.String(),
		})
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Println("success status failed")
		return
	}

	return
}

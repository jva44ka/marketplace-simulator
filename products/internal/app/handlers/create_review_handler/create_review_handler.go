package create_review_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/16/students/week-1-workshop/internal/domain/model"
	http2 "gitlab.ozon.dev/16/students/week-1-workshop/pkg/http"
	"net/http"
	"strconv"
)

type ReviewService interface {
	AddReview(ctx context.Context, review model.Review) (model.Review, error)
}

type CreateReviewHandler struct {
	reviewService ReviewService
}

func NewCreateReviewHandler(reviewService ReviewService) *CreateReviewHandler {
	return &CreateReviewHandler{reviewService: reviewService}
}

func (h CreateReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	var createReviewRequest CreateReviewRequest

	if err := json.NewDecoder(r.Body).Decode(&createReviewRequest); err != nil {
		if err = http2.ErrorResponse(w, http.StatusBadRequest, err.Error()); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	userUUID, err := uuid.Parse(createReviewRequest.UserID)
	if err != nil {
		if err = http2.ErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	review := model.Review{
		Sku:     model.Sku(createReviewRequest.Sku),
		Comment: createReviewRequest.Comment,
		UserID:  userUUID,
	}

	newReview, err := h.reviewService.AddReview(r.Context(), review)
	if err != nil {
		if err = http2.ErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {
			return
		}

		return
	}

	createReviewResponse := CreateReviewResponse{
		ID:      uint64(newReview.ID),
		Sku:     uint64(newReview.Sku),
		Comment: newReview.Comment,
		UserID:  newReview.UserID.String(),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&createReviewResponse); err != nil {
		return
	}

	return
}

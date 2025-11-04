package get_product_by_sku_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jva44ka/ozon-simulator-go/internal/domain/model"
	http2 "github.com/jva44ka/ozon-simulator-go/pkg/http"
	"net/http"
	"strconv"
)

type ProductService interface {
	GetProductsBySku(ctx context.Context, sku uint64) ([]model.Product, error)
}

type GetProductsBySkuHandler struct {
	ProductService ProductService
}

func NewGetProductsBySkuHandler(ProductService ProductService) *GetProductsBySkuHandler {
	return &GetProductsBySkuHandler{ProductService: ProductService}
}

func (h GetProductsBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	Products, err := h.ProductService.GetProductsBySku(r.Context(), sku)
	if err != nil {
		if err = http2.ErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {

			return
		}

		return
	}

	response := GetProductsResponse{Products: make([]GetProductsProductResponse, 0, len(Products))}
	for _, Product := range Products {
		response.Products = append(response.Products, GetProductsProductResponse{
			ID:      uint64(Product.ID),
			Sku:     uint64(Product.Sku),
			Comment: Product.Comment,
			UserID:  Product.UserID.String(),
		})
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Println("success status failed")
		return
	}

	return
}

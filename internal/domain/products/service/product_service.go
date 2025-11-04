package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/16/students/week-1-workshop/internal/domain/model"
	"net/http"
)

const HeaderXApiKey = "X-API-KEY"

type ProductService struct {
	client  http.Client
	token   string
	address string
}

func NewProductService(client http.Client,
	token string,
	address string) *ProductService {
	return &ProductService{
		client:  client,
		token:   token,
		address: address,
	}
}

func (s *ProductService) GetProductBySku(ctx context.Context, sku model.Sku) (*model.Product, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/product/%d", s.address, sku),
		http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	request.Header.Add(HeaderXApiKey, s.token)

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, model.ErrProductNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("http query failed")
	}

	product := &model.Product{}
	if err := json.NewDecoder(response.Body).Decode(product); err != nil {
		return nil, err
	}

	return product, nil
}

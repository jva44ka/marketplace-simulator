package get_product_by_sku_handler

type GetProductProductResponse struct {
	Sku   uint64  `json:"sku"`
	Price float64 `json:"price"`
	Name  string  `json:"name"`
}

type GetProductsResponse struct {
	Product GetProductProductResponse `json:"product"`
}

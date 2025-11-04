package get_reviews_by_sku_handler

type GetReviewsReviewResponse struct {
	ID      uint64 `json:"id"`
	Sku     uint64 `json:"sku"`
	Comment string `json:"comment"`
	UserID  string `json:"user_id"`
}

type GetReviewsResponse struct {
	Reviews []GetReviewsReviewResponse `json:"reviews"`
}

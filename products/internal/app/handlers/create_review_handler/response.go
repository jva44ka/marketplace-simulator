package create_review_handler

type CreateReviewResponse struct {
	ID      uint64 `json:"id"`
	Sku     uint64 `json:"sku"`
	Comment string `json:"comment"`
	UserID  string `json:"user_id"`
}

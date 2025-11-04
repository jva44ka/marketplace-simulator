package create_review_handler

type CreateReviewRequest struct {
	Sku     uint64 `json:"sku"`
	Comment string `json:"comment"`
	UserID  string `json:"user_id"`
}

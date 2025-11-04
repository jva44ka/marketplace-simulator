package model

import "github.com/google/uuid"

type ReviewID uint64

type Review struct {
	ID      ReviewID
	Sku     Sku
	Comment string
	UserID  uuid.UUID
}

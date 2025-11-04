package model

type Sku uint64

type Product struct {
	Sku   Sku
	Price float64
	Name  string
}

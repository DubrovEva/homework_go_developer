package models

type Product struct {
	SKU   int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}

// Products - map[sku]count
type Products = map[int64]uint32

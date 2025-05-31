package models

type Product struct {
	SKU   int64  `json:"sku"`
	Count uint32 `json:"count"`
}

package models

type Stock struct {
	SKU      int64  `json:"sku"`
	Count    uint32 `json:"total_count"`
	Reserved uint32 `json:"reserved"`
}

type StocksMap = map[int64]Stock

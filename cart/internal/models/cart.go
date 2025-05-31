package models

type Cart struct {
	Items      []*Product `json:"items"`
	TotalPrice int64      `json:"total_price"`
}

package api

type AddProductRequest struct {
	Count uint32 `json:"count" valid:"required"`
}

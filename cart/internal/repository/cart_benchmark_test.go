package repository

import (
	"testing"

	"route256/cart/internal/models"
)

var (
	defaultProduct = &models.Product{
		SKU:   123,
		Name:  "test",
		Price: 100,
		Count: 1,
	}
)

func BenchmarkCart_AddProduct(b *testing.B) {
	cart := NewCart()

	for i := 0; i < b.N; i++ {
		cart.AddProduct(defaultUserID, defaultProduct, 1)
	}
}

func BenchmarkCart_DeleteProduct(b *testing.B) {
	cart := NewCart()

	for i := 0; i < b.N; i++ {
		cart.DeleteProduct(defaultUserID, defaultProduct)
	}
}

func BenchmarkCart_Cart(b *testing.B) {
	cart := NewCart()
	cart.AddProduct(defaultUserID, defaultProduct, 1)

	for i := 0; i < b.N; i++ {
		cart.Cart(defaultUserID)
	}
}

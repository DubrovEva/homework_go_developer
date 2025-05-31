package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	"route256/cart/internal/models"
)

const (
	defaultUserID = 111
)

func TestCart_AddProduct(t *testing.T) {
	count := uint32(3)

	testProduct := &models.Product{
		SKU:   123,
		Name:  "test",
		Price: 100,
		Count: 1,
	}
	cartRepo := NewCart()

	// try to add product to empty cart
	cartRepo.AddProduct(defaultUserID, testProduct, count)

	require.Equal(t, 1, len(cartRepo.Cart(defaultUserID)))
	require.Equal(t, count, cartRepo.Cart(defaultUserID)[testProduct.SKU])

	// try to add product to non-empty cart
	cartRepo.AddProduct(defaultUserID, testProduct, count)

	require.Equal(t, 1, len(cartRepo.Cart(defaultUserID)))
	require.Equal(t, count*2, cartRepo.Cart(defaultUserID)[testProduct.SKU])
}

func TestCart_DeleteProduct(t *testing.T) {
	count := uint32(3)

	testProduct := &models.Product{
		SKU:   123,
		Name:  "test",
		Price: 100,
		Count: 1,
	}
	cartRepo := NewCart()

	// try to delete product from empty cart
	cartRepo.DeleteProduct(defaultUserID, testProduct)

	require.Equal(t, 0, len(cartRepo.Cart(defaultUserID)))

	// try to delete product from non-empty cart
	cartRepo.AddProduct(defaultUserID, testProduct, count)

	require.Equal(t, 1, len(cartRepo.Cart(defaultUserID)))
	require.Equal(t, count, cartRepo.Cart(defaultUserID)[testProduct.SKU])

	cartRepo.DeleteProduct(defaultUserID, testProduct)

	require.Equal(t, 0, len(cartRepo.Cart(defaultUserID)))
}

func TestCart_Clear(t *testing.T) {
	cartRepo := NewCart()

	require.Equal(t, 0, len(cartRepo.Cart(defaultUserID)))

	// try to clear empty cart
	cartRepo.Clear(defaultUserID)

	require.Equal(t, 0, len(cartRepo.Cart(defaultUserID)))

	// try to clear non-empty cart
	cartRepo.AddProduct(defaultUserID, &models.Product{Count: 1}, 1)
	require.Equal(t, 1, len(cartRepo.Cart(defaultUserID)))

	cartRepo.Clear(defaultUserID)

	require.Equal(t, 0, len(cartRepo.Cart(defaultUserID)))
}

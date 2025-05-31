package e2e

import (
	"route256/cart/pkg/cart"
	"route256/cart/tests"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CartSuite struct {
	suite.Suite

	client *cart.Client
}

func (s *CartSuite) SetupSuite() {
	s.client = cart.NewClient(tests.CartServiceAddress)
}

func TestCartSuite(t *testing.T) {
	if !tests.IsContainerRun() {
		t.Skip()

		return
	}

	suite.Run(t, new(CartSuite))
}

func (s *CartSuite) TestCart() {
	// add test product
	s.NoError(s.client.AddProduct(tests.DefaultUserID, tests.DefaultSKU))

	// check cart isn't empty
	userCart, err := s.client.Cart(tests.DefaultUserID)
	s.NoError(err)
	s.Equal(1, len(userCart.Items))
	s.Equal(tests.DefaultSKU, userCart.Items[0].SKU)
	s.Equal(int64(userCart.Items[0].Price), userCart.TotalPrice)

	// add another product
	s.NoError(s.client.AddProduct(tests.DefaultUserID, tests.DefaultSKU))

	// check total price is correct
	userCart, err = s.client.Cart(tests.DefaultUserID)
	s.NoError(err)
	s.Equal(1, len(userCart.Items))
	s.Equal(tests.DefaultSKU, userCart.Items[0].SKU)
	s.Equal(int64(userCart.Items[0].Price*2), userCart.TotalPrice)

	s.NoError(s.client.ClearCart(tests.DefaultUserID))
}

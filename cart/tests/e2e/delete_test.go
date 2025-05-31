package e2e

import (
	"route256/cart/tests"
	"testing"

	"github.com/stretchr/testify/suite"

	cart "route256/cart/pkg/cart"
)

type DeleteSuite struct {
	suite.Suite

	client *cart.Client
}

func (s *DeleteSuite) SetupSuite() {
	s.client = cart.NewClient(tests.CartServiceAddress)
}

func TestDeleteSuite(t *testing.T) {
	if !tests.IsContainerRun() {
		t.Skip()

		return
	}

	suite.Run(t, new(DeleteSuite))
}

func (s *DeleteSuite) TestDelete() {
	// add test product
	s.NoError(s.client.AddProduct(tests.DefaultUserID, tests.DefaultSKU))

	// check cart isn't empty
	userCart, err := s.client.Cart(tests.DefaultUserID)
	s.NoError(err)
	s.Equal(1, len(userCart.Items))
	s.Equal(tests.DefaultSKU, userCart.Items[0].SKU)
	s.Equal(int64(userCart.Items[0].Price), userCart.TotalPrice)

	// delete test product
	s.NoError(s.client.DeleteProduct(tests.DefaultUserID, tests.DefaultSKU))

	// check cart is empty
	userCart, err = s.client.Cart(tests.DefaultUserID)
	s.NoError(err)
	s.Empty(userCart)
}

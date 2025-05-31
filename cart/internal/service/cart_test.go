package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"route256/cart/internal/models"
	"route256/cart/internal/service/mocked"
)

const (
	defaultUserID = 111
)

type CartSuite struct {
	suite.Suite

	repository     *mocked.RepositoryMock
	productService *mocked.ProductsMock
	lomsService    *mocked.LomsMock

	testProduct *models.Product
	testError   error

	cart *Cart
}

func (s *CartSuite) SetupTest() {
	mc := minimock.NewController(s.T())

	s.productService = mocked.NewProductsMock(mc)
	s.repository = mocked.NewRepositoryMock(mc)
	s.lomsService = mocked.NewLomsMock(mc)

	s.cart = NewCart(s.repository, s.productService, s.lomsService)

	s.testProduct = &models.Product{
		SKU:   123,
		Price: 100,
		Name:  "test",
		Count: 1,
	}
	s.testError = errors.New("test error")
}

func TestCartSuite(t *testing.T) {
	suite.Run(t, new(CartSuite))
}

func (s *CartSuite) TestAddProduct_Ok() {
	s.productService.ProductBySkuMock.Expect(context.Background(), s.testProduct.SKU).Return(s.testProduct, nil)
	s.repository.AddProductMock.Expect(defaultUserID, s.testProduct, 1)
	s.lomsService.IsProductAvailableMock.Expect(context.Background(), s.testProduct.SKU, 1).Return(true, nil)

	err := s.cart.AddProduct(context.Background(), defaultUserID, 123, 1)
	s.Nil(err)
}

func (s *CartSuite) TestAddProduct_Err() {
	s.productService.ProductBySkuMock.Expect(context.Background(), s.testProduct.SKU).Return(nil, s.testError)

	err := s.cart.AddProduct(context.Background(), defaultUserID, 123, 1)
	s.Error(err, "can't get product by sku: test error")
}

func (s *CartSuite) TestDeleteProduct_Ok() {
	s.productService.ProductBySkuMock.Expect(context.Background(), s.testProduct.SKU).Return(s.testProduct, nil)
	s.repository.DeleteProductMock.Expect(defaultUserID, s.testProduct)

	err := s.cart.DeleteProduct(context.Background(), defaultUserID, s.testProduct.SKU)
	s.Nil(err)
}

func (s *CartSuite) TestDeleteProduct_Err() {
	s.productService.ProductBySkuMock.Expect(context.Background(), s.testProduct.SKU).Return(nil, s.testError)

	err := s.cart.DeleteProduct(context.Background(), defaultUserID, s.testProduct.SKU)
	s.Error(err, "can't get product by sku: test error")
}

func (s *CartSuite) TestCart_Ok() {
	count := uint32(3)
	secondProduct := &models.Product{
		SKU:   456,
		Price: 200,
		Name:  "test2",
		Count: 1,
	}

	s.repository.CartMock.Expect(defaultUserID).Return(map[int64]uint32{
		s.testProduct.SKU: count,
		secondProduct.SKU: count,
	})

	s.productService.ProductBySkuMock.When(minimock.AnyContext, s.testProduct.SKU).Then(s.testProduct, nil)
	s.productService.ProductBySkuMock.When(minimock.AnyContext, secondProduct.SKU).Then(secondProduct, nil)

	cart, err := s.cart.Cart(context.Background(), defaultUserID)
	s.Nil(err)
	s.Len(cart.Items, 2)

	s.Equal(s.testProduct, cart.Items[0])
	s.Equal(count, cart.Items[0].Count)

	s.Equal(secondProduct, cart.Items[1])
	s.Equal(count, cart.Items[1].Count)

	s.Equal(int64(count*s.testProduct.Price+count*secondProduct.Price), cart.TotalPrice)
	s.True(cart.Items[0].SKU < cart.Items[1].SKU)
}

func (s *CartSuite) TestCart_Err() {
	s.repository.CartMock.Expect(defaultUserID).Return(map[int64]uint32{s.testProduct.SKU: 1})
	s.productService.ProductBySkuMock.Expect(minimock.AnyContext, s.testProduct.SKU).Return(nil, s.testError)

	_, err := s.cart.Cart(context.Background(), defaultUserID)
	s.Error(err, "can't get product by sku: test error")
}

func (s *CartSuite) TestClear() {
	s.repository.ClearMock.Expect(defaultUserID).Return(true)

	wasCleared := s.cart.ClearCart(defaultUserID)
	s.Equal(true, wasCleared)
}

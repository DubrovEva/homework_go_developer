package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"route256/cart/internal/api/mocked"
	"route256/cart/internal/models"
)

const (
	defaultUserID = 111
)

type RouterSuite struct {
	suite.Suite

	cart *mocked.CartServiceMock

	router *Router

	testProduct *models.Product
	testError   error
}

func (s *RouterSuite) SetupSuite() {
	mc := minimock.NewController(s.T())

	s.cart = mocked.NewCartServiceMock(mc)

	s.router = NewRouter("localhost", "8080", s.cart)

	s.testProduct = &models.Product{
		SKU:   123,
		Price: 100,
		Name:  "test",
		Count: 1,
	}
	s.testError = errors.New("test error")
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

// addProduct

func (s *RouterSuite) TestAddProduct_Ok() {
	addProductRequest := AddProductRequest{Count: 2}
	inputBody, err := json.Marshal(addProductRequest)
	s.NoError(err)

	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(inputBody))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))

	w := httptest.NewRecorder()

	s.cart.AddProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
		addProductRequest.Count,
	).Return(nil)

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusOK)
}

func (s *RouterSuite) TestAddProduct_EmptySKU() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"sku is empty"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_ParsingSKUErr() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("sku", "invalid")
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"can't parse sku: strconv.Atoi: parsing \"invalid\": invalid syntax"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_InvalidSKU() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("sku", "0")
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"sku must be more than zero"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_EmptyUserID() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"user_id is empty"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_ParsingUserIDErr() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", "invalid")
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"can't parse user_id: strconv.Atoi: parsing \"invalid\": invalid syntax"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_InvalidUserID() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", "0")
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"user_id must be more than zero"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_EmptyBody() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"EOF"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_InvalidBody() {
	addProductRequest := AddProductRequest{Count: 0}
	inputBody, err := json.Marshal(addProductRequest)
	s.NoError(err)

	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(inputBody))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"count: non zero value required"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_InternalErr() {
	addProductRequest := AddProductRequest{Count: 2}
	inputBody, err := json.Marshal(addProductRequest)
	s.NoError(err)

	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(inputBody))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))

	w := httptest.NewRecorder()

	s.cart.AddProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
		addProductRequest.Count,
	).Return(s.testError)

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusInternalServerError)
	s.Equal(w.Body.String(), `{"message":"test error"}`+"\n")
}

func (s *RouterSuite) TestAddProduct_ProductNotFound() {
	addProductRequest := AddProductRequest{Count: 2}
	inputBody, err := json.Marshal(addProductRequest)
	s.NoError(err)

	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(inputBody))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))

	w := httptest.NewRecorder()

	s.cart.AddProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
		addProductRequest.Count,
	).Return(models.ErrProductNotFound)

	s.router.addProduct(w, req)

	s.Equal(w.Code, http.StatusPreconditionFailed)
	s.Equal(w.Body.String(), `{"message":"product not found"}`+"\n")
}

// deleteProduct

func (s *RouterSuite) TestDeleteProduct_Ok() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.cart.DeleteProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
	).Return(nil)

	s.router.deleteProduct(w, req)

	s.Equal(w.Code, http.StatusNoContent)
}

func (s *RouterSuite) TestDeleteProduct_InvalidSku() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", "0")
	w := httptest.NewRecorder()

	s.router.deleteProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"sku must be more than zero"}`+"\n")
}

func (s *RouterSuite) TestDeleteProduct_InvalidUserID() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", "0")
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))
	w := httptest.NewRecorder()

	s.router.deleteProduct(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"user_id must be more than zero"}`+"\n")
}

func (s *RouterSuite) TestDeleteProduct_InternalErr() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))

	w := httptest.NewRecorder()

	s.cart.DeleteProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
	).Return(s.testError)

	s.router.deleteProduct(w, req)

	s.Equal(w.Code, http.StatusInternalServerError)
	s.Equal(w.Body.String(), `{"message":"test error"}`+"\n")
}

func (s *RouterSuite) TestDeleteProduct_ProductNotFound() {
	path := fmt.Sprintf("/user/%d/cart/%d", defaultUserID, s.testProduct.SKU)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	req.SetPathValue("sku", fmt.Sprintf("%d", s.testProduct.SKU))

	w := httptest.NewRecorder()

	s.cart.DeleteProductMock.Expect(
		context.Background(),
		defaultUserID,
		s.testProduct.SKU,
	).Return(models.ErrProductNotFound)

	s.router.deleteProduct(w, req)

	s.Equal(w.Code, http.StatusPreconditionFailed)
	s.Equal(w.Body.String(), `{"message":"product not found"}`+"\n")
}

// cart

func (s *RouterSuite) TestCart_Ok() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	w := httptest.NewRecorder()

	s.cart.CartMock.Expect(
		context.Background(),
		defaultUserID,
	).Return(&models.Cart{
		Items: []*models.Product{
			s.testProduct,
		},
		TotalPrice: int64(s.testProduct.Price),
	}, nil)

	s.router.cart(w, req)

	s.Equal(w.Code, http.StatusOK)
	s.Equal(w.Body.String(), `{"items":[{"sku":123,"name":"test","price":100,"count":1}],"total_price":100}`+"\n")
}

func (s *RouterSuite) TestCart_Empty() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	w := httptest.NewRecorder()

	s.cart.CartMock.Expect(
		context.Background(),
		defaultUserID,
	).Return(&models.Cart{}, nil)

	s.router.cart(w, req)

	s.Equal(w.Code, http.StatusNotFound)
}

func (s *RouterSuite) TestCart_InvalidUserID() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", "0")
	w := httptest.NewRecorder()

	s.router.cart(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"user_id must be more than zero"}`+"\n")
}

func (s *RouterSuite) TestCart_InternalErr() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	w := httptest.NewRecorder()

	s.cart.CartMock.Expect(
		context.Background(),
		defaultUserID,
	).Return(nil, s.testError)

	s.router.cart(w, req)

	s.Equal(w.Code, http.StatusInternalServerError)
	s.Equal(w.Body.String(), `{"message":"test error"}`+"\n")
}

// clearCart

func (s *RouterSuite) TestClearCart_Ok() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", fmt.Sprintf("%d", defaultUserID))
	w := httptest.NewRecorder()

	s.cart.ClearCartMock.Expect(defaultUserID).Return(true)

	s.router.clearCart(w, req)

	s.Equal(w.Code, http.StatusNoContent)
}

func (s *RouterSuite) TestClearCart_InvalidUserID() {
	path := fmt.Sprintf("/user/%d/cart", defaultUserID)

	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewReader(nil))
	req.SetPathValue("user_id", "0")
	w := httptest.NewRecorder()

	s.router.clearCart(w, req)

	s.Equal(w.Code, http.StatusBadRequest)
	s.Equal(w.Body.String(), `{"message":"user_id must be more than zero"}`+"\n")
}

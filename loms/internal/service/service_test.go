package service

//import (
//	"context"
//	"errors"
//	"testing"
//
//	"github.com/gojuno/minimock/v3"
//	"github.com/stretchr/testify/suite"
//
//	"route256/loms/internal/service/mocked"
//	loms "route256/loms/pkg/api/loms/v1"
//)
//
//type ServiceSuite struct {
//	suite.Suite
//
//	ordersMock *mocked.OrdersMock
//	stocksMock *mocked.StocksMock
//	dbMock     *mocked.DBMock
//
//	service   *Service
//	testError error
//}
//
//func (s *ServiceSuite) SetupTest() {
//	mc := minimock.NewController(s.T())
//
//	s.ordersMock = mocked.NewOrdersMock(mc)
//	s.stocksMock = mocked.NewStocksMock(mc)
//	s.dbMock = mocked.NewDBMock(mc)
//
//	s.service = NewService(s.dbMock, s.ordersMock, s.stocksMock)
//
//	s.testError = errors.New("test error")
//}
//
//func TestServiceSuite(t *testing.T) {
//	suite.Run(t, new(ServiceSuite))
//}
//
//// CreateOrder
//
//func (s *ServiceSuite) TestCreateOrder_Ok() {
//	items := []*loms.Item{{SkuId: 1, Count: 2}}
//	orderID := int64(123)
//
//	s.ordersMock.CreateMock.Return(orderID, nil)
//	s.stocksMock.ReserveMock.Return(nil)
//	s.ordersMock.AwaitingPaymentMock.Return(nil)
//
//	result, err := s.service.CreateOrder(context.Background(), 1, items)
//
//	s.Nil(err)
//	s.Equal(orderID, result)
//}
//
//func (s *ServiceSuite) TestCreateOrder_CreateErr() {
//	items := []*loms.Item{{SkuId: 1, Count: 2}}
//	orderID := int64(123)
//
//	s.ordersMock.CreateMock.Return(orderID, s.testError)
//
//	result, err := s.service.CreateOrder(context.Background(), 1, items)
//
//	s.Error(err)
//	s.Zero(result)
//}
//
//func (s *ServiceSuite) TestCreateOrder_ReserveErr() {
//	items := []*loms.Item{{SkuId: 1, Count: 2}}
//	orderID := int64(123)
//
//	s.ordersMock.CreateMock.Return(orderID, nil)
//	s.stocksMock.ReserveMock.Return(s.testError)
//	s.ordersMock.FailMock.Return(nil)
//
//	result, err := s.service.CreateOrder(context.Background(), 1, items)
//
//	s.Error(err)
//	s.Zero(result)
//}
//
//func (s *ServiceSuite) TestCreateOrder_FailErr() {
//	items := []*loms.Item{{SkuId: 1, Count: 2}}
//	orderID := int64(123)
//
//	s.ordersMock.CreateMock.Return(orderID, nil)
//	s.stocksMock.ReserveMock.Return(s.testError)
//	s.ordersMock.FailMock.Return(s.testError)
//
//	result, err := s.service.CreateOrder(context.Background(), 1, items)
//
//	s.Error(err)
//	s.Zero(result)
//}
//
//func (s *ServiceSuite) TestCreateOrder_AwaitingPaymentErr() {
//	items := []*loms.Item{{SkuId: 1, Count: 2}}
//	orderID := int64(123)
//
//	s.ordersMock.CreateMock.Return(orderID, nil)
//	s.stocksMock.ReserveMock.Return(nil)
//	s.ordersMock.AwaitingPaymentMock.Return(s.testError)
//
//	result, err := s.service.CreateOrder(context.Background(), 1, items)
//
//	s.Error(err)
//	s.Zero(result)
//}
//
//// OrderInfo
//
//func (s *ServiceSuite) TestOrderInfo_Ok() {
//	order := &loms.Order{Id: 123}
//	s.ordersMock.OrderMock.Return(order, nil)
//
//	result, err := s.service.OrderInfo(context.Background(), 123)
//
//	s.Nil(err)
//	s.Equal(order, result)
//}
//
//func (s *ServiceSuite) TestOrderInfo_OrderErr() {
//	s.ordersMock.OrderMock.Return(nil, s.testError)
//
//	result, err := s.service.OrderInfo(context.Background(), 123)
//
//	s.Error(err)
//	s.Nil(result)
//}
//
//// PayOrder
//
//func (s *ServiceSuite) TestPayOrder_Ok() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{
//		{SkuId: 1, Count: 2},
//		{SkuId: 2, Count: 3},
//	}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.SellMock.Return(nil)
//	s.ordersMock.PayMock.Return(nil)
//
//	err := s.service.PayOrder(context.Background(), 123)
//
//	s.Nil(err)
//}
//
//func (s *ServiceSuite) TestPayOrder_OrderErr() {
//	s.ordersMock.OrderMock.Return(nil, s.testError)
//
//	err := s.service.PayOrder(context.Background(), 123)
//	s.Error(err)
//}
//
//func (s *ServiceSuite) TestPayOrder_SellErr() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{
//		{SkuId: 1, Count: 2},
//		{SkuId: 2, Count: 3},
//	}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.SellMock.Return(s.testError)
//
//	err := s.service.PayOrder(context.Background(), 123)
//	s.Error(err)
//}
//
//func (s *ServiceSuite) TestPayOrder_PayErr() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{
//		{SkuId: 1, Count: 2},
//		{SkuId: 2, Count: 3},
//	}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.SellMock.Return(nil)
//	s.ordersMock.PayMock.Return(s.testError)
//
//	err := s.service.PayOrder(context.Background(), 123)
//	s.Error(err)
//}
//
//// CancelOrder
//
//func (s *ServiceSuite) TestCancelOrder_Ok() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{{SkuId: 1, Count: 2}}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.ReleaseMock.Return(nil)
//	s.ordersMock.CancelMock.Return(nil)
//
//	err := s.service.CancelOrder(context.Background(), 123)
//
//	s.Nil(err)
//}
//
//func (s *ServiceSuite) TestCancelOrder_ReleaseErr() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{{SkuId: 1, Count: 2}}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.ReleaseMock.Return(s.testError)
//
//	err := s.service.CancelOrder(context.Background(), 123)
//
//	s.Error(err)
//}
//
//func (s *ServiceSuite) TestCancelOrder_OrderErr() {
//	s.ordersMock.OrderMock.Return(nil, s.testError)
//
//	err := s.service.CancelOrder(context.Background(), 123)
//
//	s.Error(err)
//}
//
//func (s *ServiceSuite) TestCancelOrder_StocksErr() {
//	order := &loms.Order{Id: 123, Items: []*loms.Item{{SkuId: 1, Count: 2}}}
//	s.ordersMock.OrderMock.Return(order, nil)
//	s.stocksMock.ReleaseMock.Return(nil)
//	s.ordersMock.CancelMock.Return(s.testError)
//
//	err := s.service.CancelOrder(context.Background(), 123)
//
//	s.Error(err)
//}
//
//// StocksInfo
//
//func (s *ServiceSuite) TestStocksInfo_Ok() {
//	s.stocksMock.AvailableStocksMock.Return(uint32(10), nil)
//
//	stocks, err := s.service.StocksInfo(context.Background(), 1)
//
//	s.Nil(err)
//	s.Equal(uint32(10), stocks)
//}
//
//func (s *ServiceSuite) TestStocksInfo_AvailableStocksErr() {
//	s.stocksMock.AvailableStocksMock.Return(uint32(0), s.testError)
//
//	stocks, err := s.service.StocksInfo(context.Background(), 1)
//
//	s.Error(err)
//	s.Zero(stocks)
//}

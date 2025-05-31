package repository

//import (
//	"testing"
//
//	"route256/loms/internal/models"
//	loms "route256/loms/pkg/api/loms/v1"
//
//	"github.com/stretchr/testify/require"
//)
//
//// Create
//
//func TestOrders_CreateAndRetrieve(t *testing.T) {
//	repo := NewOrders()
//	userID := int64(123)
//	items := []*loms.Item{
//		{SkuId: 1, Count: 2},
//	}
//
//	orderID := repo.Create(userID, items)
//	require.NotZero(t, orderID, "order ID should not be zero")
//
//	order, err := repo.Order(orderID)
//	require.NoError(t, err, "retrieving the order should not return an error")
//	require.Equal(t, orderID, order.Id)
//	require.Equal(t, userID, order.UserId)
//	require.Equal(t, loms.OrderStatus_NEW, order.Status)
//	require.Equal(t, items, order.Items)
//}
//
//// Order
//
//func TestOrders_Order_NotFound(t *testing.T) {
//	repo := NewOrders()
//
//	_, err := repo.Order(999)
//	require.ErrorIs(t, err, models.ErrOrderNotFound, "should return ErrOrderNotFound")
//}
//
//func TestOrders_Order_Ok(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//
//	order, err := repo.Order(orderID)
//	require.NoError(t, err, "Order should not return an error")
//	require.Equal(t, orderID, order.Id)
//}
//
//// AwaitingPayment
//
//func TestOrders_AwaitingPayment(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//
//	err := repo.AwaitingPayment(orderID)
//	require.NoError(t, err, "AwaitingPayment should not return an error")
//
//	order, _ := repo.Order(orderID)
//	require.Equal(t, loms.OrderStatus_AWAITING_PAYMENT, order.Status)
//}
//
//func TestOrders_AwaitingPayment_InvalidStatus(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//
//	err := repo.AwaitingPayment(orderID)
//	require.ErrorIs(t, err, models.ErrOrderIsNotNew)
//}
//
//func TestOrders_AwaitingPayment_NotFound(t *testing.T) {
//	repo := NewOrders()
//
//	err := repo.AwaitingPayment(999)
//	require.ErrorIs(t, err, models.ErrOrderNotFound)
//}
//
//// Pay
//
//func TestOrders_Pay(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//
//	err := repo.Pay(orderID)
//	require.NoError(t, err, "Pay should not return an error")
//
//	order, _ := repo.Order(orderID)
//	require.Equal(t, loms.OrderStatus_PAYED, order.Status)
//}
//
//func TestOrders_Pay_InvalidStatus(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//
//	err := repo.Pay(orderID)
//	require.ErrorIs(t, err, models.ErrOrderIsNotAwaitingPayment)
//}
//
//func TestOrders_Pay_AlreadyPaid(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//	require.NoError(t, repo.Pay(orderID))
//
//	err := repo.Pay(orderID)
//	require.NoError(t, err, "Pay should not return an error")
//}
//
//func TestOrders_Pay_NotFound(t *testing.T) {
//	repo := NewOrders()
//
//	err := repo.Pay(999)
//	require.ErrorIs(t, err, models.ErrOrderNotFound)
//}
//
//// Cancel
//
//func TestOrders_Cancel(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//
//	err := repo.Cancel(orderID)
//	require.NoError(t, err, "Cancel should not return an error")
//
//	order, _ := repo.Order(orderID)
//	require.Equal(t, loms.OrderStatus_CANCELLED, order.Status)
//}
//
//func TestOrders_Cancel_NotFound(t *testing.T) {
//	repo := NewOrders()
//
//	err := repo.Cancel(999)
//	require.ErrorIs(t, err, models.ErrOrderNotFound)
//}
//
//func TestOrders_Cancel_AlreadyCancelled(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//	require.NoError(t, repo.Cancel(orderID))
//
//	err := repo.Cancel(orderID)
//	require.NoError(t, err, "Cancel should not return an error")
//}
//
//func TestOrders_Cancel_AlreadyPaid(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//	require.NoError(t, repo.Pay(orderID))
//
//	err := repo.Cancel(orderID)
//	require.ErrorIs(t, err, models.ErrOrderIsAlreadyPayed)
//}
//
//// Fail
//
//func TestOrders_Fail_Ok(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//
//	err := repo.Fail(orderID)
//	require.NoError(t, err, "Fail should not return an error")
//
//	order, _ := repo.Order(orderID)
//	require.Equal(t, loms.OrderStatus_FAILED, order.Status)
//}
//
//func TestOrders_Fail_AlreadyPaid(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//	require.NoError(t, repo.Pay(orderID))
//
//	err := repo.Fail(orderID)
//	require.ErrorIs(t, err, models.ErrOrderIsAlreadyPayed)
//}
//
//func TestOrders_Fail_AlreadyCancelled(t *testing.T) {
//	repo := NewOrders()
//	orderID := repo.Create(123, nil)
//	require.NoError(t, repo.AwaitingPayment(orderID))
//	require.NoError(t, repo.Cancel(orderID))
//
//	err := repo.Fail(orderID)
//	require.NoError(t, err, "Fail should not return an error")
//}
//
//func TestOrders_Fail_NotFound(t *testing.T) {
//	repo := NewOrders()
//
//	err := repo.Fail(999)
//	require.ErrorIs(t, err, models.ErrOrderNotFound)
//}

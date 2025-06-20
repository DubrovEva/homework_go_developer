// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package generated

import (
	"context"
)

type Querier interface {
	AddOrder(ctx context.Context, userID int64) (*Order, error)
	AddOrderItem(ctx context.Context, arg *AddOrderItemParams) (*OrdersItem, error)
	CancelOrder(ctx context.Context, id int64) error
	FailOrder(ctx context.Context, id int64) error
	Order(ctx context.Context, id int64) (*Order, error)
	OrderItems(ctx context.Context, orderID int64) ([]*OrdersItem, error)
	PayOrder(ctx context.Context, id int64) error
	SetAwaitingPayment(ctx context.Context, id int64) error
	Stock(ctx context.Context, sku int64) (*Stock, error)
	UpdateStock(ctx context.Context, arg *UpdateStockParams) error
}

var _ Querier = (*Queries)(nil)

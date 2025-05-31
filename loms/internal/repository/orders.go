package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"route256/loms/internal/repository/generated"

	"github.com/jackc/pgx/v5/pgxpool"
	"route256/loms/internal/models"
	loms "route256/loms/pkg/api/loms/v1"
)

type Orders struct {
	pool *pgxpool.Pool
}

func NewOrders(pool *pgxpool.Pool) *Orders {
	return &Orders{
		pool: pool,
	}
}

func (o *Orders) Create(ctx context.Context, userID int64, items []*loms.Item) (int64, error) {
	repository := generated.New(o.pool)

	order, err := repository.AddOrder(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("can't add order: %w", err)
	}

	for _, item := range items {
		params := &generated.AddOrderItemParams{
			OrderID: order.ID,
			Sku:     item.SkuId,
			Count:   int32(item.Count),
		}

		_, err := repository.AddOrderItem(ctx, params)
		if err != nil {
			return 0, fmt.Errorf("can't add order item: %w", err)
		}
	}

	return order.ID, nil
}

func (o *Orders) Order(ctx context.Context, id int64) (*loms.Order, error) {
	repository := generated.New(o.pool)

	order, err := repository.Order(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}

		return nil, fmt.Errorf("can't get order: %w", err)
	}

	lomsOrder := loms.Order{
		Id:     order.ID,
		Status: loms.OrderStatus(order.Status),
		UserId: order.UserID,
	}

	list, err := repository.OrderItems(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("can't get order items: %w", err)
	}

	lomsItems := make([]*loms.Item, 0, len(list))
	for _, item := range list {
		if item.Count < 0 {
			log.Printf("negative count of item in order: %d, sku: %d", item.Count, item.Sku)

			return nil, models.ErrNegativeCount
		}

		lomsItems = append(lomsItems, &loms.Item{
			SkuId: item.Sku,
			Count: uint32(item.Count),
		})
	}
	lomsOrder.Items = lomsItems

	return &lomsOrder, nil
}

func (o *Orders) Pay(ctx context.Context, id int64) error {
	repository := generated.New(o.pool)

	order, err := repository.Order(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrOrderNotFound
		}

		return fmt.Errorf("can't get order: %w", err)
	}

	if order.Status == int32(loms.OrderStatus_PAYED) {
		return nil
	}

	if order.Status != int32(loms.OrderStatus_AWAITING_PAYMENT) {
		return models.ErrOrderIsNotAwaitingPayment
	}

	if err := repository.PayOrder(ctx, id); err != nil {
		return fmt.Errorf("can't update order status: %w", err)
	}

	return nil
}

func (o *Orders) AwaitingPayment(ctx context.Context, id int64) error {
	repository := generated.New(o.pool)

	order, err := repository.Order(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrOrderNotFound
		}

		return fmt.Errorf("can't get order: %w", err)
	}

	if order.Status != int32(loms.OrderStatus_NEW) {
		return models.ErrOrderIsNotNew
	}

	if err := repository.SetAwaitingPayment(ctx, id); err != nil {
		return fmt.Errorf("can't update order status: %w", err)
	}

	return nil
}

func (o *Orders) Cancel(ctx context.Context, id int64) error {
	repository := generated.New(o.pool)

	order, err := repository.Order(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrOrderNotFound
		}

		return fmt.Errorf("can't get order: %w", err)
	}

	if order.Status == int32(loms.OrderStatus_CANCELLED) {
		return nil
	}

	if order.Status == int32(loms.OrderStatus_PAYED) {
		return models.ErrOrderIsAlreadyPayed
	}

	if err := repository.CancelOrder(ctx, id); err != nil {
		return fmt.Errorf("can't update order status: %w", err)
	}

	return nil
}

func (o *Orders) Fail(ctx context.Context, id int64) error {
	repository := generated.New(o.pool)

	order, err := repository.Order(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrOrderNotFound
		}

		return fmt.Errorf("can't get order: %w", err)
	}

	if order.Status == int32(loms.OrderStatus_CANCELLED) {
		return nil
	}

	if order.Status == int32(loms.OrderStatus_PAYED) {
		return models.ErrOrderIsAlreadyPayed
	}

	if err := repository.FailOrder(ctx, id); err != nil {
		return fmt.Errorf("can't update order status: %w", err)
	}

	return nil
}

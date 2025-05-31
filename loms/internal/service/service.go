package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"route256/loms/internal/kafka"
	"route256/loms/internal/models"
	loms "route256/loms/pkg/api/loms/v1"
)

type Service struct {
	db DB

	orders Orders
	stocks Stocks

	kafkaProducer *kafka.Producer
}

//go:generate minimock -i route256/loms/internal/service.DB -o mocked/db_mock.go -n DBMock -p mocked

type DB interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

//go:generate minimock -i route256/loms/internal/service.Orders -o mocked/orders_mock.go -n OrdersMock -p mocked

type Orders interface {
	Create(ctx context.Context, userID int64, items []*loms.Item) (int64, error)
	Order(ctx context.Context, id int64) (*loms.Order, error)
	Pay(ctx context.Context, id int64) error
	AwaitingPayment(ctx context.Context, id int64) error
	Cancel(ctx context.Context, id int64) error
	Fail(ctx context.Context, id int64) error
}

//go:generate minimock -i route256/loms/internal/service.Stocks -o mocked/stocks_mock.go -n StocksMock -p mocked

type Stocks interface {
	Reserve(ctx context.Context, items []*loms.Item) error
	Release(ctx context.Context, items []*loms.Item) error
	Sell(ctx context.Context, items []*loms.Item) error
	AvailableStocks(ctx context.Context, sku int64) (uint32, error)
}

func NewService(db DB, orders Orders, stocks Stocks, kafkaProducer *kafka.Producer) *Service {
	return &Service{
		db:            db,
		orders:        orders,
		stocks:        stocks,
		kafkaProducer: kafkaProducer,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, items []*loms.Item) (int64, error) {
	var orderID int64

	txErr := pgx.BeginTxFunc(ctx, s.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var err error

		orderID, err = s.orders.Create(ctx, userID, items)
		if err != nil {
			return fmt.Errorf("can't create order: %w", err)
		}

		s.sendOrderEvent(ctx, orderID, "created", loms.OrderStatus_NEW.String(), userID, nil)

		if err = s.stocks.Reserve(ctx, items); err != nil {
			failedErr := s.orders.Fail(ctx, orderID)
			if failedErr != nil {
				return fmt.Errorf("can't fail order: %w, reservation error: %w", failedErr, err)
			}

			s.sendOrderEvent(ctx, orderID, "failed", loms.OrderStatus_FAILED.String(), userID, map[string]string{
				"reason": "reservation_failed",
			})

			return fmt.Errorf("%w, reservation error: %w", models.ErrOrderFailed, err)
		}

		if err = s.orders.AwaitingPayment(ctx, orderID); err != nil {
			return fmt.Errorf("can't update order status: %w", err)
		}

		s.sendOrderEvent(ctx, orderID, "awaiting_payment", loms.OrderStatus_AWAITING_PAYMENT.String(), userID, nil)

		return nil
	})
	if txErr != nil {
		return 0, fmt.Errorf("can't create order: %w", txErr)
	}

	return orderID, nil
}

func (s *Service) OrderInfo(ctx context.Context, orderID int64) (*loms.Order, error) {
	order, err := s.orders.Order(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("can't get order: %w", err)
	}

	return order, nil
}

func (s *Service) PayOrder(ctx context.Context, orderID int64) error {
	var userID int64

	txErr := pgx.BeginTxFunc(ctx, s.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		order, err := s.orders.Order(ctx, orderID)
		if err != nil {
			return fmt.Errorf("can't get order: %w", err)
		}

		userID = order.UserId

		if err = s.stocks.Sell(ctx, order.Items); err != nil {
			return fmt.Errorf("can't pay order: %w", err)
		}

		if err = s.orders.Pay(ctx, orderID); err != nil {
			return fmt.Errorf("can't update order status: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return fmt.Errorf("can't pay order: %w", txErr)
	}

	s.sendOrderEvent(ctx, orderID, "paid", loms.OrderStatus_PAYED.String(), userID, nil)

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, orderID int64) error {
	var userID int64

	txErr := pgx.BeginTxFunc(ctx, s.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		order, err := s.orders.Order(ctx, orderID)
		if err != nil {
			return fmt.Errorf("can't get order: %w", err)
		}

		userID = order.UserId

		if err = s.stocks.Release(ctx, order.Items); err != nil {
			return fmt.Errorf("can't release order items: %w", err)
		}

		if err = s.orders.Cancel(ctx, orderID); err != nil {
			return fmt.Errorf("can't update order status: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return fmt.Errorf("can't cancel order: %w", txErr)
	}

	s.sendOrderEvent(ctx, orderID, "cancelled", loms.OrderStatus_CANCELLED.String(), userID, nil)

	return nil
}

func (s *Service) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	stocksCount, err := s.stocks.AvailableStocks(ctx, sku)
	if err != nil {
		return 0, fmt.Errorf("can't get stocks info: %w", err)
	}

	return stocksCount, nil
}

func (s *Service) sendOrderEvent(ctx context.Context, orderID int64, event, status string, userID int64, additional any) {
	if s.kafkaProducer == nil {
		return
	}

	orderEvent := kafka.OrderEvent{
		OrderID:    orderID,
		Event:      event,
		Status:     status,
		Timestamp:  time.Now(),
		UserID:     userID,
		Additional: additional,
	}

	err := s.kafkaProducer.SendOrderEvent(ctx, orderEvent)
	if err != nil {
		fmt.Printf("Failed to send order event to Kafka: %v\n", err)
	}
}

package api

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/loms/internal/models"
	loms "route256/loms/pkg/api/loms/v1"
)

var _ loms.LomsServer = (*Loms)(nil)

type Loms struct {
	loms.UnimplementedLomsServer

	service Service
}

type Service interface {
	CreateOrder(ctx context.Context, userID int64, items []*loms.Item) (int64, error)
	OrderInfo(ctx context.Context, id int64) (*loms.Order, error)
	PayOrder(ctx context.Context, id int64) error
	CancelOrder(ctx context.Context, id int64) error
	StocksInfo(ctx context.Context, sku int64) (uint32, error)
}

func NewLoms(service Service) *Loms {
	return &Loms{
		service: service,
	}
}

func (s *Loms) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	orderID, err := s.service.CreateOrder(ctx, req.UserId, req.Items)
	if err != nil {
		if errors.Is(err, models.ErrOrderFailed) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &loms.CreateOrderResponse{
		OrderId: orderID,
	}

	return response, nil
}

func (s *Loms) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	order, err := s.service.OrderInfo(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &loms.OrderInfoResponse{
		Order: order,
	}

	return response, nil
}

func (s *Loms) PayOrder(ctx context.Context, req *loms.PayOrderRequest) (*loms.PayOrderResponse, error) {
	if err := s.service.PayOrder(ctx, req.OrderId); err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, models.ErrOrderIsNotAwaitingPayment) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &loms.PayOrderResponse{}, nil
}

func (s *Loms) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*loms.CancelOrderResponse, error) {
	if err := s.service.CancelOrder(ctx, req.OrderId); err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, models.ErrOrderIsAlreadyPayed) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &loms.CancelOrderResponse{}, nil
}

func (s *Loms) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	stocksCount, err := s.service.StocksInfo(ctx, req.Sku)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &loms.StocksInfoResponse{
		Count: stocksCount,
	}

	return response, nil
}

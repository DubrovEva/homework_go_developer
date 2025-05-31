package loms_service

import (
	"context"
	"fmt"
	"time"

	"route256/cart/internal/logging"
	"route256/cart/internal/metrics"
	"route256/cart/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	loms "route256/cart/pkg/loms/api/v1"
)

type Client struct {
	client loms.LomsClient
	tracer trace.Tracer
}

func NewClient(port int64, host string) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %v", err)
	}

	logging.Info("Connected to gRPC server", zap.String("target", target))

	client := loms.NewLomsClient(conn)

	return &Client{
		client: client,
		tracer: tracing.Tracer("cart-loms-service"),
	}, nil
}

func (c *Client) CreateOrder(ctx context.Context, items []*loms.Item) (int64, error) {
	startTime := time.Now()

	ctx, span := c.tracer.Start(ctx, "LomsService.CreateOrder")
	defer span.End()

	traceID, spanID := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(attribute.Int("items.count", len(items)))

	logging.Info("Creating order in LOMS",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	request := &loms.CreateOrderRequest{
		Items: items,
	}

	response, err := c.client.CreateOrder(ctx, request)
	duration := time.Since(startTime)

	statusCode := 0
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = int(st.Code())
		}
	}

	metrics.ObserveExternalRequest("loms-service", "CreateOrder", statusCode, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logging.Error("Failed to create order in LOMS",
			zap.Error(err),
			zap.Int("items_count", len(items)),
			zap.String("trace_id", traceID),
		)

		return 0, fmt.Errorf("failed to create order: %v", err)
	}

	logging.Info("Successfully created order in LOMS",
		zap.Int64("order_id", response.OrderId),
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	return response.OrderId, nil
}

func (c *Client) IsProductAvailable(ctx context.Context, sku int64, count uint32) (bool, error) {
	startTime := time.Now()

	ctx, span := c.tracer.Start(ctx, "LomsService.IsProductAvailable")
	defer span.End()

	traceID, spanID := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(
		attribute.Int64("product.sku", sku),
		attribute.Int64("product.count", int64(count)),
	)

	logging.Info("Checking product availability in LOMS",
		zap.Int64("sku", sku),
		zap.Uint32("count", count),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	request := &loms.StocksInfoRequest{
		Sku: sku,
	}

	response, err := c.client.StocksInfo(ctx, request)
	duration := time.Since(startTime)

	statusCode := 0
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = int(st.Code())
		}
	}

	metrics.ObserveExternalRequest("loms-service", "StocksInfo", statusCode, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logging.Error("Failed to get stocks info from LOMS",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("trace_id", traceID),
		)

		return false, fmt.Errorf("failed to get stocks info: %v", err)
	}

	isAvailable := response.Count >= count

	logging.Info("Product availability check completed",
		zap.Int64("sku", sku),
		zap.Uint32("requested_count", count),
		zap.Uint32("available_count", response.Count),
		zap.Bool("is_available", isAvailable),
		zap.String("trace_id", traceID),
	)

	return isAvailable, nil
}

package product_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/logging"
	"route256/cart/internal/metrics"
	"route256/cart/internal/models"
	"route256/cart/internal/tracing"
	"route256/cart/pkg/ratelimiter"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	HeaderXApiKey = "X-API-KEY"

	defaultRetryCount = 3
)

type Client struct {
	client  http.Client
	limiter *ratelimiter.RateLimiter

	token   string
	address string

	retryCount uint
	tracer     trace.Tracer
}

func NewClient(
	ctx context.Context, token, schema, host, port string, retryCount uint, rps int,
) *Client {
	if retryCount == 0 {
		retryCount = defaultRetryCount
	}

	client := http.Client{
		Transport: NewMiddleware(http.DefaultTransport, retryCount),
	}

	address := fmt.Sprintf("%s://%s:%s", schema, host, port)
	limiter := ratelimiter.NewRateLimiter(ctx, rps)

	return &Client{
		client:     client,
		limiter:    limiter,
		token:      token,
		address:    address,
		retryCount: retryCount,
		tracer:     tracing.Tracer("cart-product-service"),
	}
}

func (c *Client) ProductBySku(ctx context.Context, sku int64) (*models.Product, error) {
	startTime := time.Now()

	ctx, span := c.tracer.Start(ctx, "ProductService.ProductBySku")
	defer span.End()

	span.SetAttributes(attribute.Int64("product.sku", sku))

	traceID, spanID := tracing.ExtractTraceInfoFromContext(ctx)

	logging.Info("Requesting product information",
		zap.Int64("sku", sku),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	if err := c.limiter.Allow(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logging.Error("Rate limiter error",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("trace_id", traceID),
		)
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	url := fmt.Sprintf("%s/product/%d", c.address, sku)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		http.NoBody,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logging.Error("Failed to create request",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("url", url),
			zap.String("trace_id", traceID),
		)
		return nil, fmt.Errorf("can't create request to get product: %w", err)
	}

	request.Header.Add(HeaderXApiKey, c.token)

	// Inject trace context into the request headers
	propagator := propagation.TraceContext{}
	propagator.Inject(ctx, propagation.HeaderCarrier(request.Header))

	response, err := c.client.Do(request)
	duration := time.Since(startTime)

	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
	}
	metrics.ObserveExternalRequest("product-service", "ProductBySku", statusCode, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logging.Error("Failed to get product",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("url", url),
			zap.String("trace_id", traceID),
		)
		return nil, fmt.Errorf("can't get product: %w", err)
	}

	defer response.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", response.StatusCode))

	if response.StatusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "product not found")
		logging.Warn("Product not found",
			zap.Int64("sku", sku),
			zap.String("trace_id", traceID),
		)
		return nil, models.ErrProductNotFound
	}

	if response.StatusCode != http.StatusOK {
		err := errors.New("http query failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logging.Error("HTTP query failed",
			zap.Int("status_code", response.StatusCode),
			zap.Int64("sku", sku),
			zap.String("url", url),
			zap.String("trace_id", traceID),
		)
		return nil, err
	}

	product := &models.Product{}
	if err := json.NewDecoder(response.Body).Decode(product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logging.Error("Failed to decode product",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("trace_id", traceID),
		)
		return nil, fmt.Errorf("can't decode product: %w", err)
	}

	logging.Info("Successfully retrieved product",
		zap.Int64("sku", sku),
		zap.Uint32("price", product.Price),
		zap.String("name", product.Name),
		zap.String("trace_id", traceID),
	)

	return product, nil
}

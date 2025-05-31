package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"route256/loms/internal/logging"
	"route256/loms/internal/metrics"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/generated"
	"route256/loms/internal/tracing"

	loms "route256/loms/pkg/api/loms/v1"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Stocks struct {
	pool *pgxpool.Pool
}

func NewStocks(pool *pgxpool.Pool) *Stocks {
	return &Stocks{
		pool: pool,
	}
}

func (s *Stocks) Reserve(ctx context.Context, items []*loms.Item) error {
	startTime := time.Now()

	ctx, span := tracing.Tracer("loms-repository").Start(ctx, "Stocks.Reserve")
	defer span.End()

	traceID, _ := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(attribute.Int("items.count", len(items)))

	logging.Info("Reserving stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	repository := generated.New(s.pool)

	for _, item := range items {
		stockStartTime := time.Now()

		stock, err := repository.Stock(ctx, item.SkuId)

		stockDuration := time.Since(stockStartTime)
		metrics.ObserveDatabaseOperation("select", err, stockDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			if errors.Is(err, pgx.ErrNoRows) {
				logging.Error("Product not found",
					zap.Int64("sku", item.SkuId),
					zap.String("trace_id", traceID),
				)
				return models.ErrProductNotFound
			}

			logging.Error("Failed to get stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't get stock: %w", err)
		}

		if stock.TotalCount < int32(item.Count) {
			err := models.ErrNotEnoughStockLeft
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Not enough stock left",
				zap.Int64("sku", item.SkuId),
				zap.Int32("available", stock.TotalCount),
				zap.Uint32("requested", item.Count),
				zap.String("trace_id", traceID),
			)
			return err
		}

		updateStartTime := time.Now()

		err = repository.UpdateStock(ctx, &generated.UpdateStockParams{
			Sku:        item.SkuId,
			Reserved:   stock.Reserved + int32(item.Count),
			TotalCount: stock.TotalCount,
		})

		updateDuration := time.Since(updateStartTime)
		metrics.ObserveDatabaseOperation("update", err, updateDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Failed to reserve stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.Uint32("count", item.Count),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't reserve stock: %w", err)
		}

		logging.Info("Reserved stock",
			zap.Int64("sku", item.SkuId),
			zap.Uint32("count", item.Count),
			zap.String("trace_id", traceID),
		)
	}

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("transaction", nil, duration)

	logging.Info("Successfully reserved all stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	return nil
}

func (s *Stocks) Release(ctx context.Context, items []*loms.Item) error {
	startTime := time.Now()

	ctx, span := tracing.Tracer("loms-repository").Start(ctx, "Stocks.Release")
	defer span.End()

	traceID, _ := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(attribute.Int("items.count", len(items)))

	logging.Info("Releasing stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	repository := generated.New(s.pool)

	for _, item := range items {
		stockStartTime := time.Now()

		stock, err := repository.Stock(ctx, item.SkuId)

		stockDuration := time.Since(stockStartTime)
		metrics.ObserveDatabaseOperation("select", err, stockDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			if errors.Is(err, pgx.ErrNoRows) {
				logging.Error("Product not found",
					zap.Int64("sku", item.SkuId),
					zap.String("trace_id", traceID),
				)
				return models.ErrProductNotFound
			}

			logging.Error("Failed to get stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't get stock: %w", err)
		}

		if stock.Reserved < int32(item.Count) {
			err := models.ErrNotEnoughStockReserved
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Not enough stock reserved",
				zap.Int64("sku", item.SkuId),
				zap.Int32("reserved", stock.Reserved),
				zap.Uint32("requested", item.Count),
				zap.String("trace_id", traceID),
			)
			return err
		}

		updateStartTime := time.Now()

		err = repository.UpdateStock(ctx, &generated.UpdateStockParams{
			Sku:        item.SkuId,
			Reserved:   stock.Reserved - int32(item.Count),
			TotalCount: stock.TotalCount,
		})

		updateDuration := time.Since(updateStartTime)
		metrics.ObserveDatabaseOperation("update", err, updateDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Failed to release stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.Uint32("count", item.Count),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't release stock: %w", err)
		}

		logging.Info("Released stock",
			zap.Int64("sku", item.SkuId),
			zap.Uint32("count", item.Count),
			zap.String("trace_id", traceID),
		)
	}

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("transaction", nil, duration)

	logging.Info("Successfully released all stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	return nil
}

func (s *Stocks) Sell(ctx context.Context, items []*loms.Item) error {
	startTime := time.Now()

	ctx, span := tracing.Tracer("loms-repository").Start(ctx, "Stocks.Sell")
	defer span.End()

	traceID, _ := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(attribute.Int("items.count", len(items)))

	logging.Info("Selling stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	repository := generated.New(s.pool)

	for _, item := range items {
		stockStartTime := time.Now()

		stock, err := repository.Stock(ctx, item.SkuId)

		stockDuration := time.Since(stockStartTime)
		metrics.ObserveDatabaseOperation("select", err, stockDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			if errors.Is(err, pgx.ErrNoRows) {
				logging.Error("Product not found",
					zap.Int64("sku", item.SkuId),
					zap.String("trace_id", traceID),
				)
				return models.ErrProductNotFound
			}

			logging.Error("Failed to get stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't get stock: %w", err)
		}

		if stock.TotalCount < int32(item.Count) {
			err := models.ErrNotEnoughStockLeft
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Not enough stock left",
				zap.Int64("sku", item.SkuId),
				zap.Int32("available", stock.TotalCount),
				zap.Uint32("requested", item.Count),
				zap.String("trace_id", traceID),
			)
			return err
		}

		if stock.Reserved < int32(item.Count) {
			err := models.ErrNotEnoughStockReserved
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Not enough stock reserved",
				zap.Int64("sku", item.SkuId),
				zap.Int32("reserved", stock.Reserved),
				zap.Uint32("requested", item.Count),
				zap.String("trace_id", traceID),
			)
			return err
		}

		updateStartTime := time.Now()

		err = repository.UpdateStock(ctx, &generated.UpdateStockParams{
			Sku:        item.SkuId,
			Reserved:   stock.Reserved - int32(item.Count),
			TotalCount: stock.TotalCount - int32(item.Count),
		})

		updateDuration := time.Since(updateStartTime)
		metrics.ObserveDatabaseOperation("update", err, updateDuration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logging.Error("Failed to sell stock",
				zap.Error(err),
				zap.Int64("sku", item.SkuId),
				zap.Uint32("count", item.Count),
				zap.String("trace_id", traceID),
			)
			return fmt.Errorf("can't sell stock: %w", err)
		}

		logging.Info("Sold stock",
			zap.Int64("sku", item.SkuId),
			zap.Uint32("count", item.Count),
			zap.String("trace_id", traceID),
		)
	}

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("transaction", nil, duration)

	logging.Info("Successfully sold all stocks",
		zap.Int("items_count", len(items)),
		zap.String("trace_id", traceID),
	)

	return nil
}

func (s *Stocks) AvailableStocks(ctx context.Context, sku int64) (uint32, error) {
	startTime := time.Now()

	ctx, span := tracing.Tracer("loms-repository").Start(ctx, "Stocks.AvailableStocks")
	defer span.End()

	traceID, _ := tracing.ExtractTraceInfoFromContext(ctx)

	span.SetAttributes(attribute.Int64("product.sku", sku))

	logging.Info("Getting available stocks",
		zap.Int64("sku", sku),
		zap.String("trace_id", traceID),
	)

	repository := generated.New(s.pool)

	stock, err := repository.Stock(ctx, sku)

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("select", err, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, pgx.ErrNoRows) {
			logging.Error("Product not found",
				zap.Int64("sku", sku),
				zap.String("trace_id", traceID),
			)
			return 0, models.ErrProductNotFound
		}

		logging.Error("Failed to get stock",
			zap.Error(err),
			zap.Int64("sku", sku),
			zap.String("trace_id", traceID),
		)
		return 0, fmt.Errorf("can't get stock: %w", err)
	}

	stocksLeft := stock.TotalCount - stock.Reserved
	if stocksLeft < 0 {
		err := models.ErrNegativeStocksCountLeft
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logging.Error("Negative stocks count left",
			zap.Int64("sku", sku),
			zap.Int32("total", stock.TotalCount),
			zap.Int32("reserved", stock.Reserved),
			zap.Int32("left", stocksLeft),
			zap.String("trace_id", traceID),
		)
		return 0, err
	}

	logging.Info("Available stocks retrieved",
		zap.Int64("sku", sku),
		zap.Int32("available", stocksLeft),
		zap.String("trace_id", traceID),
	)

	return uint32(stocksLeft), nil
}

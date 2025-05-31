package repository

import (
	"route256/cart/internal/logging"
	"route256/cart/internal/metrics"
	"route256/cart/internal/models"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Cart struct {
	carts map[int64]models.Products

	mx sync.RWMutex
}

func NewCart() *Cart {
	return &Cart{
		carts: make(map[int64]models.Products),
	}
}

func (c *Cart) AddProduct(userID int64, product *models.Product, count uint32) {
	startTime := time.Now()

	c.mx.Lock()
	defer c.mx.Unlock()

	if _, exist := c.carts[userID]; !exist {
		c.carts[userID] = make(models.Products, 1)
	}

	c.carts[userID][product.SKU] += count

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("insert", nil, duration)

	logging.Info("Added product to cart",
		zap.Int64("user_id", userID),
		zap.Int64("sku", product.SKU),
		zap.Uint32("count", count),
		zap.Uint32("total_count", c.carts[userID][product.SKU]),
	)
}

func (c *Cart) DeleteProduct(userID int64, product *models.Product) {
	startTime := time.Now()

	c.mx.Lock()
	defer c.mx.Unlock()

	if _, exist := c.carts[userID]; !exist {
		logging.Warn("Attempted to delete product from non-existent cart",
			zap.Int64("user_id", userID),
			zap.Int64("sku", product.SKU),
		)
		return
	}

	delete(c.carts[userID], product.SKU)

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("delete", nil, duration)

	logging.Info("Deleted product from cart",
		zap.Int64("user_id", userID),
		zap.Int64("sku", product.SKU),
	)
}

func (c *Cart) Cart(userID int64) models.Products {
	startTime := time.Now()

	c.mx.RLock()
	defer c.mx.RUnlock()

	if _, exist := c.carts[userID]; !exist {
		logging.Info("Cart not found",
			zap.Int64("user_id", userID),
		)
		return nil
	}

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("select", nil, duration)

	logging.Info("Retrieved cart",
		zap.Int64("user_id", userID),
		zap.Int("items_count", len(c.carts[userID])),
	)

	return c.carts[userID]
}

func (c *Cart) Clear(userID int64) bool {
	startTime := time.Now()

	c.mx.Lock()
	defer c.mx.Unlock()

	if _, exist := c.carts[userID]; !exist {
		logging.Warn("Attempted to clear non-existent cart",
			zap.Int64("user_id", userID),
		)
		return false
	}

	delete(c.carts, userID)

	duration := time.Since(startTime)
	metrics.ObserveDatabaseOperation("delete", nil, duration)

	logging.Info("Cleared cart",
		zap.Int64("user_id", userID),
	)

	return true
}

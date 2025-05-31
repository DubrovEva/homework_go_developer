package service

import (
	"context"
	"fmt"
	"route256/cart/pkg/errgroup"
	"slices"
	"sync"

	loms "route256/cart/pkg/loms/api/v1"

	"route256/cart/internal/models"
)

type Cart struct {
	repository      Repository
	productsService Products
	lomsService     Loms
}

//go:generate minimock -i route256/cart/internal/service.Products -o mocked/products_mock.go -n ProductsMock -p mocked

type Products interface {
	ProductBySku(ctx context.Context, sku int64) (*models.Product, error)
}

//go:generate minimock -i route256/cart/internal/service.Loms -o mocked/loms_mock.go -n LomsMock -p mocked

type Loms interface {
	CreateOrder(ctx context.Context, items []*loms.Item) (int64, error)
	IsProductAvailable(ctx context.Context, sku int64, count uint32) (bool, error)
}

//go:generate minimock -i route256/cart/internal/service.Repository -o mocked/repository_mock.go -n RepositoryMock -p mocked

type Repository interface {
	Clear(userID int64) bool
	Cart(userID int64) models.Products
	AddProduct(userID int64, product *models.Product, count uint32)
	DeleteProduct(userID int64, product *models.Product)
}

func NewCart(repository Repository, productsService Products, lomsService Loms) *Cart {
	return &Cart{
		repository:      repository,
		productsService: productsService,
		lomsService:     lomsService,
	}
}

func (c *Cart) AddProduct(ctx context.Context, userID, sku int64, count uint32) error {
	product, err := c.productsService.ProductBySku(ctx, sku)
	if err != nil {
		return fmt.Errorf("can't get product by sku: %w", err)
	}

	isAvailable, err := c.lomsService.IsProductAvailable(ctx, sku, count)
	if err != nil {
		return fmt.Errorf("can't check product availability: %w", err)
	}

	if !isAvailable {
		return fmt.Errorf("product is not available")
	}

	c.repository.AddProduct(userID, product, count)

	return nil
}

func (c *Cart) DeleteProduct(ctx context.Context, userID, sku int64) error {
	product, err := c.productsService.ProductBySku(ctx, sku)
	if err != nil {
		return fmt.Errorf("can't get product by sku: %w", err)
	}

	c.repository.DeleteProduct(userID, product)

	return nil
}

func (c *Cart) Cart(ctx context.Context, userID int64) (*models.Cart, error) {
	products := c.repository.Cart(userID)
	cart := models.Cart{
		Items: make([]*models.Product, 0, len(products)),
	}

	var (
		mx     sync.Mutex
		result = make([]*models.Product, 0, len(products))
	)

	eg, grCtx := errgroup.WithContext(ctx)

	for sku, count := range products {
		sku, count := sku, count // Копируем переменные для корректной работы в горутине (нужно для старых версий go)

		eg.Go(func() error {
			product, err := c.productsService.ProductBySku(grCtx, sku)
			if err != nil {
				return fmt.Errorf("can't get product by sku %d: %w", sku, err)
			}

			product.Count = count

			mx.Lock()
			cart.TotalPrice += int64(product.Price * count)
			result = append(result, product)
			mx.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	cart.Items = result

	slices.SortFunc(cart.Items, func(a, b *models.Product) int {
		return int(a.SKU - b.SKU)
	})

	return &cart, nil
}

func (c *Cart) ClearCart(userID int64) bool {
	return c.repository.Clear(userID)
}

func (c *Cart) Checkout(ctx context.Context, userID int64) (int64, error) {
	items, err := c.items(userID)
	if err != nil {
		return 0, fmt.Errorf("can't get items: %w", err)
	}

	return c.lomsService.CreateOrder(ctx, items)
}

func (c *Cart) items(userID int64) ([]*loms.Item, error) {
	products := c.repository.Cart(userID)

	items := make([]*loms.Item, 0, len(products))

	for sku, count := range products {
		items = append(items, &loms.Item{
			SkuId: sku,
			Count: count,
		})
	}

	return items, nil
}

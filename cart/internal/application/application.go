package application

import (
	"context"
	"fmt"

	"route256/cart/internal/api"
	"route256/cart/internal/config"
	"route256/cart/internal/repository"
	"route256/cart/internal/service"
	"route256/cart/internal/service/loms_service"
	"route256/cart/internal/service/product_service"
)

type App struct {
	config *config.Config
	router *api.Router

	cartRepo *repository.Cart

	cartService    *service.Cart
	productService *product_service.Client
	lomsService    *loms_service.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	a.initRepository()

	if err := a.initService(ctx); err != nil {
		return fmt.Errorf("can't init service: %w", err)
	}

	if err := a.initRouter(); err != nil {
		return fmt.Errorf("can't init router: %w", err)
	}

	return nil
}

func (a *App) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	a.config = cfg

	return nil
}

func (a *App) initRepository() {
	a.cartRepo = repository.NewCart()

}

func (a *App) initService(ctx context.Context) error {
	a.productService = product_service.NewClient(
		ctx,
		a.config.Products.Token,
		a.config.Products.Schema,
		a.config.Products.Host,
		a.config.Products.Port,
		a.config.Products.RetryCount,
		a.config.Products.RPS,
	)
	lomsService, err := loms_service.NewClient(a.config.Loms.Port, a.config.Loms.Host)
	if err != nil {
		return fmt.Errorf("can't init loms service: %w", err)
	}

	a.lomsService = lomsService
	a.cartService = service.NewCart(a.cartRepo, a.productService, a.lomsService)

	return nil
}

func (a *App) initRouter() error {
	a.router = api.NewRouter(a.config.Service.Host, a.config.Service.Port, a.cartService)

	if err := a.router.Start(); err != nil {
		return fmt.Errorf("can't start router: %w", err)
	}

	return nil
}

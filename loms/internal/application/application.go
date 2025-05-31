package application

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"route256/loms/internal/api"
	"route256/loms/internal/config"
	"route256/loms/internal/kafka"
	"route256/loms/internal/repository"
	"route256/loms/internal/service"
)

type App struct {
	config *config.Config

	gateway *api.Gateway
	loms    *api.Loms

	service *service.Service

	stocksRepo *repository.Stocks
	ordersRepo *repository.Orders

	pool *pgxpool.Pool

	kafkaProducer *kafka.Producer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	if err := a.initRepository(ctx); err != nil {
		return fmt.Errorf("can't init repository: %w", err)
	}

	if err := a.initKafkaProducer(); err != nil {
		return fmt.Errorf("can't init kafka producer: %w", err)
	}

	a.initService()

	if err := a.initLoms(ctx); err != nil {
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

func (a *App) initRepository(ctx context.Context) error {
	var err error

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		a.config.DB.User,
		a.config.DB.Password,
		a.config.DB.Host,
		a.config.DB.Port,
		a.config.DB.DBName,
	)

	fmt.Printf("%v", a.config.DB)

	fmt.Println(connStr)

	masterConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("unable to parse config: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, masterConfig)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v\n", err)
	}

	a.stocksRepo = repository.NewStocks(pool)
	a.ordersRepo = repository.NewOrders(pool)

	a.pool = pool

	return nil
}

func (a *App) initKafkaProducer() error {
	var err error
	a.kafkaProducer, err = kafka.NewProducer(a.config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	return nil
}

func (a *App) initService() {
	a.service = service.NewService(a.pool, a.ordersRepo, a.stocksRepo, a.kafkaProducer)
}

func (a *App) initLoms(ctx context.Context) error {
	a.loms = api.NewLoms(a.service)
	a.gateway = api.NewGateway(a.loms, a.config.Service.GrpcPort, a.config.Service.HttpPort)

	if err := a.gateway.Run(ctx); err != nil {
		return fmt.Errorf("can't start router: %w", err)
	}

	return nil
}

// Close gracefully shuts down the application
func (a *App) Close() error {
	log.Println("Shutting down LOMS application...")

	if a.kafkaProducer != nil {
		if err := a.kafkaProducer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
			return err
		}
		log.Println("Kafka producer closed successfully")
	}

	if a.pool != nil {
		a.pool.Close()
		log.Println("Database connection pool closed")
	}

	return nil
}

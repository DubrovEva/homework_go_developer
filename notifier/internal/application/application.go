package application

import (
	"context"
	"fmt"
	"log"

	"route256/notifier/internal/config"
	"route256/notifier/internal/kafka"
)

type App struct {
	config        *config.Config
	kafkaConsumer *kafka.Consumer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	if err := a.initKafkaConsumer(); err != nil {
		return fmt.Errorf("can't init kafka consumer: %w", err)
	}

	log.Printf("Notifier instance %s started", a.config.Service.InstanceID)

	return a.kafkaConsumer.Start(ctx)
}

func (a *App) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	a.config = cfg

	return nil
}

func (a *App) initKafkaConsumer() error {
	var err error
	a.kafkaConsumer, err = kafka.NewConsumer(a.config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	return nil
}

func (a *App) Close() error {
	if a.kafkaConsumer != nil {
		if err := a.kafkaConsumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
			return err
		}
	}
	return nil
}

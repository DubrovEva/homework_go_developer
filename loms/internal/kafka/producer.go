package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"route256/loms/internal/config"

	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	OrderID    int64     `json:"order_id"`
	Event      string    `json:"event"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     int64     `json:"user_id,omitempty"`
	Additional any       `json:"additional,omitempty"`
}

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers),
		Topic:                  cfg.Kafka.OrderTopic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
		topic:  cfg.Kafka.OrderTopic,
	}, nil
}

func (p *Producer) SendOrderEvent(ctx context.Context, event OrderEvent) error {
	key := []byte(fmt.Sprintf("%d", event.OrderID))

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal order event: %w", err)
	}

	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)

		return err
	}

	return nil
}

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"route256/notifier/internal/config"

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

type Consumer struct {
	reader *kafka.Reader
	topic  string
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Kafka.Brokers},
		Topic:       cfg.Kafka.OrderTopic,
		GroupID:     cfg.Kafka.ConsumerGroupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	return &Consumer{
		reader: reader,
		topic:  cfg.Kafka.OrderTopic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("Starting Kafka consumer for topic: %s", c.topic)

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping consumer")

			return nil
		default:
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				log.Printf("Error fetching message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) processMessage(message kafka.Message) error {
	var event OrderEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	log.Printf("Received order event: OrderID=%d, Event=%s, Status=%s, Timestamp=%s",
		event.OrderID, event.Event, event.Status, event.Timestamp.Format(time.RFC3339))

	if event.UserID != 0 {
		log.Printf("  UserID: %d", event.UserID)
	}

	if event.Additional != nil {
		additionalJSON, err := json.Marshal(event.Additional)
		if err == nil {
			log.Printf("  Additional info: %s", string(additionalJSON))
		}
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)

		return err
	}

	return nil
}

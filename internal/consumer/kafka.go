package consumer

import (
	"context"
	"fmt"
	"gift-store/internal/config"
	"gift-store/internal/consumer/handlers"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DebeziumMessage struct {
	Payload DebeziumPayload `json:"payload"`
}

type DebeziumPayload struct {
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
	Op     string                 `json:"op"`
}

type Consumer struct {
	kafkaConsumer   *kafka.Consumer
	handlerRegistry *handlers.HandlerRegistry
	wg              sync.WaitGroup
}

func NewConsumer(cfg config.AppConfig, handlerRegistry *handlers.HandlerRegistry) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Kafka_Brokers,
		"group.id":           cfg.Kafka_Group_ID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	// Subscribe to all configured topics
	if err := c.SubscribeTopics(cfg.Kafka_Topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	return &Consumer{
		kafkaConsumer:   c,
		handlerRegistry: handlerRegistry,
	}, nil
}

// Start starts the consumer in a separate goroutine
func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	go c.consume(ctx)
}

// Wait waits for the consumer to finish
func (c *Consumer) Wait() {
	c.wg.Wait()
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.wg.Wait()
	return c.kafkaConsumer.Close()
}

// consume processes messages from Kafka
func (c *Consumer) consume(ctx context.Context) {
	defer c.wg.Done()

	// Create a ticker for polling with a small interval
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, shutting down...")
			return
		case <-ticker.C:
			// Use Poll instead of ReadMessage to make it more responsive to context cancellation
			event := c.kafkaConsumer.Poll(100) // 100ms timeout
			if event == nil {
				continue
			}

			switch e := event.(type) {
			case *kafka.Message:
				// Get the handler for this topic
				handler := c.handlerRegistry.GetHandler(*e.TopicPartition.Topic)
				if handler == nil {
					log.Printf("No handler found for topic: %s\n", *e.TopicPartition.Topic)
					continue
				}

				// Process the message with the appropriate handler
				if err := handler.Handle(ctx, e.Value); err != nil {
					log.Printf("Error processing message: %v\n", err)
					continue
				}

				// Commit the offset after successful processing
				if _, err := c.kafkaConsumer.CommitMessage(e); err != nil {
					log.Printf("Error committing offset: %v\n", err)
				}

			case kafka.Error:
				if e.Code() == kafka.ErrAllBrokersDown {
					log.Println("All brokers are down. Retrying...")
				} else {
					log.Printf("Kafka error: %v\n", e)
				}
			}
		}
	}
}

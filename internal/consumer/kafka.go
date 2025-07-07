package consumer

import (
	"context"
	"fmt"
	"gift-store/internal/config"
	"gift-store/internal/consumer/handlers"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	kafkaConsumer   *kafka.Consumer
	handlerRegistry *handlers.HandlerRegistry
	wg              sync.WaitGroup
}

func NewConsumer(cfg config.AppConfig, handlerRegistry *handlers.HandlerRegistry) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Kafka_Brokers,
		"group.id":           cfg.Kafka_Group_ID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	// Subscribe to all configured topics
	if err := consumer.SubscribeTopics(cfg.Kafka_Topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	return &Consumer{
		kafkaConsumer:   consumer,
		handlerRegistry: handlerRegistry,
	}, nil
}

// Start starts the consumer in a separate goroutine
func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	go c.consume(ctx)
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.wg.Wait()
	return c.kafkaConsumer.Close()
}

// consume processes messages from Kafka
func (c *Consumer) consume(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			c.wg.Done()
			c.Close()
			return
		default:
			event := c.kafkaConsumer.Poll(-1)
			if event == nil {
				continue
			}

			switch e := event.(type) {
			case *kafka.Message:
				handler := c.handlerRegistry.GetHandler(*e.TopicPartition.Topic)
				if handler == nil {
					log.Printf("No handler found for topic: %s\n", *e.TopicPartition.Topic)
					continue
				}

				if err := handler.Handle(ctx, e.Value); err != nil {
					log.Printf("Error processing message: %v\n", err)
					continue
				}

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

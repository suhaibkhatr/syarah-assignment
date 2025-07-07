package main

import (
	"context"
	"gift-store/internal/config"
	"gift-store/internal/consumer"
	"gift-store/internal/consumer/handlers"
	"gift-store/internal/db"
	"gift-store/internal/repo/elastic"
	"gift-store/internal/repo/mysql"
	"gift-store/internal/sink"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database source and sink
	source := db.NewSource(cfg)
	esSink := sink.NewSink(cfg)

	// Initial data migration from MySQL to Elasticsearch
	productRepo := mysql.NewProductRepo(source)
	products, err := productRepo.GetAllProducts()
	if err != nil {
		log.Fatal("Error reading products:", err)
	}

	productElasticRepo := elastic.NewProductRepo(esSink)
	for _, p := range products {
		err := productElasticRepo.InsertOrUpdate(p, "products")
		if err != nil {
			log.Printf("Error migrating product %d: %v", p.ID, err)
		}
	}

	log.Println("Initial data migration completed")

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the handler registry
	handlerRegistry := handlers.NewHandlerRegistry(
		handlers.NewProductHandler(esSink),
		// Add other handlers here as needed
	)

	// Create and start the consumer
	kafkaConsumer, err := consumer.NewConsumer(cfg, handlerRegistry)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start the consumer in a separate goroutine
	kafkaConsumer.Start(ctx)

	// Set up signal handling for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigchan
	log.Println("Shutting down...")
	
	// Cancel the context to signal the consumer to stop
	cancel()
	
	// Close the consumer
	if err := kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
	
	log.Println("Shutdown completed")
}

package main

import (
	"gift-store/internal/config"
	"gift-store/internal/consumer"
	"gift-store/internal/db"
	"gift-store/internal/sink"
	"log"
)

func main() {
	cfg := config.Load()

	source := db.NewSource(cfg)
	sink := sink.NewSink(cfg)

	// Initial migration
	products, err := source.GetAllProducts()
	if err != nil {
		log.Fatal("Error reading products:", err)
	}
	for _, p := range products {
		err := sink.InsertOrUpdate(p, "products")
		if err != nil {
			log.Fatal("Error inserting product:", err)
		}
	}

	// Start CDC consumer
	consumer.ListenAndSync(cfg, sink)
}

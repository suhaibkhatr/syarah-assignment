package sink

import (
	"context"
	"gift-store/internal/config"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
)

type ElasticSink struct {
	Es *elastic.Client
}

func NewSink(cfg config.AppConfig) *ElasticSink {
	// Add retry on failure and sniffing disabled for better connection handling
	es, err := elastic.NewClient(
		elastic.SetURL(cfg.Elastic_URL),
		elastic.SetSniff(false), // Disable sniffing in development
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckTimeoutStartup(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Check if Elasticsearch is available
	info, code, err := es.Ping(cfg.Elastic_URL).Do(context.Background())
	if err != nil {
		log.Fatalf("Elasticsearch is not available: %v", err)
	}
	if code != 200 {
		log.Fatalf("Elasticsearch returned status code: %d", code)
	}
	log.Printf("Connected to Elasticsearch (version: %s)", info.Version.Number)

	return &ElasticSink{Es: es}
}

package consumer

import (
	"encoding/json"
	"fmt"
	"gift-store/internal/config"
	"gift-store/internal/sink"
	"gift-store/internal/util"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Main Debezium message structure
type DebeziumMessage struct {
	Schema  DebeziumSchema  `json:"schema"`
	Payload DebeziumPayload `json:"payload"`
}

// Schema structure (optional to parse, but useful for understanding)
type DebeziumSchema struct {
	Type     string        `json:"type"`
	Fields   []SchemaField `json:"fields"`
	Optional bool          `json:"optional"`
	Name     string        `json:"name"`
	Version  int           `json:"version,omitempty"`
}

type SchemaField struct {
	Type       string            `json:"type"`
	Fields     []SchemaField     `json:"fields,omitempty"`
	Optional   bool              `json:"optional"`
	Field      string            `json:"field"`
	Name       string            `json:"name,omitempty"`
	Version    int               `json:"version,omitempty"`
	Default    interface{}       `json:"default,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// Payload structure
type DebeziumPayload struct {
	Before      map[string]interface{} `json:"before"`
	After       map[string]interface{} `json:"after"`
	Source      DebeziumSource         `json:"source"`
	Transaction *DebeziumTransaction   `json:"transaction"`
	Op          string                 `json:"op"`
	TsMs        *int64                 `json:"ts_ms"`
	TsUs        *int64                 `json:"ts_us"`
	TsNs        *int64                 `json:"ts_ns"`
}

// Source information
type DebeziumSource struct {
	Version   string  `json:"version"`
	Connector string  `json:"connector"`
	Name      string  `json:"name"`
	TsMs      int64   `json:"ts_ms"`
	Snapshot  *string `json:"snapshot"`
	Db        string  `json:"db"`
	Sequence  *string `json:"sequence"`
	TsUs      *int64  `json:"ts_us"`
	TsNs      *int64  `json:"ts_ns"`
	Table     *string `json:"table"`
	ServerID  int64   `json:"server_id"`
	Gtid      *string `json:"gtid"`
	File      string  `json:"file"`
	Pos       int64   `json:"pos"`
	Row       int32   `json:"row"`
	Thread    *int64  `json:"thread"`
	Query     *string `json:"query"`
}

// Transaction information
type DebeziumTransaction struct {
	ID                  string `json:"id"`
	TotalOrder          int64  `json:"total_order"`
	DataCollectionOrder int64  `json:"data_collection_order"`
}

func ListenAndSync(cfg config.AppConfig, sink *sink.ElasticSink) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka_Brokers,
		"group.id":          cfg.Kafka_Group_ID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = consumer.SubscribeTopics(cfg.Kafka_Topic, nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Listening for changes...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Println(err)
			continue
		}
		var event DebeziumMessage
		fmt.Println(msg.TopicPartition.Topic)
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(event.Payload.Source.Name, event.Payload.Source.Table)
		switch event.Payload.Op {
		case "c", "u":
			product, err := util.MapToProduct(event.Payload.After)
			if err != nil {
				log.Println(err)
				continue
			}
			err = sink.InsertOrUpdate(*product)
			if err != nil {
				log.Println(err)
			}
			log.Println("Product inserted or updated:", product.ID)
		case "d":
			product, err := util.MapToProduct(event.Payload.Before)
			if err != nil {
				log.Println(err)
				continue
			}
			err = sink.Delete(product.ID)
			if err != nil {
				log.Println(err)
			}
			log.Println("Product deleted:", product.ID)
		}
	}
}

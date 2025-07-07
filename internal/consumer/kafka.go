package consumer

import (
	"encoding/json"
	"fmt"
	"gift-store/internal/config"
	"gift-store/internal/models"
	"gift-store/internal/sink"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DebeziumPayload struct {
	Payload struct {
		Op     string         `json:"op"`
		Before models.Product `json:"before"`
		After  models.Product `json:"after"`
	} `json:"payload"`
}

func ListenAndSync(cfg config.AppConfig, sink *sink.ElasticSink) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s", cfg.Kafka.Brokers[0]),
		"group.id":          cfg.Kafka.GroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	err = consumer.Subscribe(cfg.Kafka.Topic, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			continue
		}
		var event DebeziumPayload
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			continue
		}
		switch event.Payload.Op {
		case "c", "u":
			err := sink.InsertOrUpdate(event.Payload.After)
			if err != nil {
				log.Println(err)
			}
		case "d":
			err := sink.Delete(event.Payload.Before.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

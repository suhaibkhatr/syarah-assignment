package consumer

import (
	"encoding/json"
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

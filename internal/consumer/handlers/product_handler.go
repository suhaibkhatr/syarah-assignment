package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gift-store/internal/repo/elastic"
	"gift-store/internal/sink"
	"gift-store/internal/util"
	"log"
)

type DebeziumMessage struct {
	Payload struct {
		Before map[string]interface{} `json:"before"`
		After  map[string]interface{} `json:"after"`
		Source map[string]interface{} `json:"source"`
		Op     string                 `json:"op"`
	} `json:"payload"`
}

type ProductHandler struct {
	sink *sink.ElasticSink
	repo *elastic.ProductRepo
}

func NewProductHandler(sink *sink.ElasticSink) *ProductHandler {
	return &ProductHandler{
		sink: sink,
		repo: elastic.NewProductRepo(sink),
	}
}

func (h *ProductHandler) CanHandle(topic string) bool {
	return topic == "gift_store.gift_store.products"
}
func (h *ProductHandler) Handle(ctx context.Context, message []byte) error {
	var msg DebeziumMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return fmt.Errorf("error unmarshaling message: %v", err)
	}

	switch msg.Payload.Op {
	case "c": // Create
		return h.handleCreate(msg.Payload.After)
	case "u": // Update
		return h.handleUpdate(msg.Payload.After)
	case "d": // Delete
		return h.handleDelete(msg.Payload.Before)
	default:
		log.Printf("Unhandled operation: %s", msg.Payload.Op)
	}

	return nil
}

func (h *ProductHandler) handleCreate(tmp map[string]interface{}) error {
	product, err := util.MapToProduct(tmp)
	if err != nil {
		return fmt.Errorf("error mapping product: %v", err)
	}
	return h.repo.InsertOrUpdate(*product, "products")
}

func (h *ProductHandler) handleUpdate(tmp map[string]interface{}) error {
	return h.handleCreate(tmp)
}

func (h *ProductHandler) handleDelete(tmp map[string]interface{}) error {
	product, err := util.MapToProduct(tmp)
	if err != nil {
		return fmt.Errorf("error mapping product: %v", err)
	}
	return h.repo.Delete(product.ID, "products")
}

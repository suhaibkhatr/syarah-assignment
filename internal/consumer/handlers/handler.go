package handlers

import "context"

// MessageHandler defines the interface for processing Kafka messages
type MessageHandler interface {
	CanHandle(topic string) bool
	Handle(ctx context.Context, message []byte) error
}

// HandlerRegistry maintains a list of available handlers
type HandlerRegistry struct {
	handlers []MessageHandler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry(handlers ...MessageHandler) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: handlers,
	}
}

// GetHandler returns the appropriate handler for the given topic
func (r *HandlerRegistry) GetHandler(topic string) MessageHandler {
	for _, h := range r.handlers {
		if h.CanHandle(topic) {
			return h
		}
	}
	return nil
}

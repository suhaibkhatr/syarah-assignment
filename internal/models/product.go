package models

import (
	"time"
)

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Category     string    `json:"category"`
	GiftCategory string    `json:"gift_category"`
	AgeGroup     string    `json:"age_group"`
	Brand        string    `json:"brand"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

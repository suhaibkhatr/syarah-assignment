package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gift-store/internal/models"
	"math/big"
	"strconv"
	"time"
)

func MapToProduct(dataMap map[string]interface{}) (*models.Product, error) {
	if dataMap == nil {
		return nil, nil
	}

	// Handle Debezium decimal (Base64 encoded price)
	if priceVal, exists := dataMap["price"]; exists {
		if priceStr, ok := priceVal.(string); ok {
			// Try to decode as Base64 first
			if decoded, err := base64.StdEncoding.DecodeString(priceStr); err == nil {
				// Convert bytes to big.Int, then to float64
				bigInt := new(big.Int).SetBytes(decoded)
				scale := 100.0 // Scale is 2 (from your schema: scale:2)
				floatVal := new(big.Float).SetInt(bigInt)
				floatVal.Quo(floatVal, big.NewFloat(scale))
				result, _ := floatVal.Float64()
				dataMap["price"] = result
			} else {
				// Try to parse as regular float string
				if parsed, err := strconv.ParseFloat(priceStr, 64); err == nil {
					dataMap["price"] = parsed
				}
			}
		}
	}

	// Handle timestamps - convert from string to time.Time
	for _, field := range []string{"created_at", "updated_at"} {
		if timeVal, exists := dataMap[field]; exists {
			if timeStr, ok := timeVal.(string); ok && timeStr != "" {
				// Try different time formats
				formats := []string{
					time.RFC3339,
					time.RFC3339Nano,
					"2006-01-02T15:04:05Z",
					"2006-01-02T15:04:05.000Z",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05.000000Z",
				}

				var parsedTime time.Time
				var err error
				for _, format := range formats {
					parsedTime, err = time.Parse(format, timeStr)
					if err == nil {
						break
					}
				}

				if err == nil {
					dataMap[field] = parsedTime
				}
			}
		}
	}

	// Handle is_available - convert to bool
	if availableVal, exists := dataMap["is_available"]; exists {
		switch v := availableVal.(type) {
		case int:
			dataMap["is_available"] = v != 0
		case int16:
			dataMap["is_available"] = v != 0
		case int32:
			dataMap["is_available"] = v != 0
		case int64:
			dataMap["is_available"] = v != 0
		case float64:
			dataMap["is_available"] = v != 0
		case string:
			if parsed, err := strconv.ParseBool(v); err == nil {
				dataMap["is_available"] = parsed
			} else if parsed, err := strconv.Atoi(v); err == nil {
				dataMap["is_available"] = parsed != 0
			}
		}
	}

	// Convert map to JSON then to Product struct
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map to JSON: %v", err)
	}

	var product models.Product
	if err := json.Unmarshal(jsonData, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Product: %v", err)
	}

	return &product, nil
}

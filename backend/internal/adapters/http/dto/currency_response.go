package dto

import (
	"time"
)

// CurrencyGetRateResponse represents the response body for currency conversion
type CurrencyGetRateResponse struct {
	Timestamp time.Time          `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

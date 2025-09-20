package dto

import (
	"time"
)

type CurrencyRateModel struct {
	Rate     float64 `json:"rate" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}
type CurrencyRatesModel struct {
	Rates     []CurrencyRateModel `json:"rates" binding:"required"`
	Timestamp time.Time           `json:"timestamp" binding:"required" format:"date-time"`
}

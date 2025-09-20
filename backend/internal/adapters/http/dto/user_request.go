package dto

type UpdatePreferredCurrencyRequest struct {
	Currency string `json:"currency" binding:"required"`
}

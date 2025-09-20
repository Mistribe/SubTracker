package dto

import (
	"time"
)

type CreateProviderRequest struct {
	Id             *string            `json:"id,omitempty"`
	Name           string             `json:"name" binding:"required"`
	Description    *string            `json:"description,omitempty"`
	IconUrl        *string            `json:"icon_url,omitempty"`
	Url            *string            `json:"url,omitempty"`
	PricingPageUrl *string            `json:"pricing_page_url,omitempty"`
	Labels         []string           `json:"labels" binding:"required"`
	Owner          EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt      *time.Time         `json:"created_at,omitempty" format:"date-time"`
}

type UpdateProviderRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    *string    `json:"description,omitempty"`
	IconUrl        *string    `json:"icon_url,omitempty"`
	Url            *string    `json:"url,omitempty"`
	PricingPageUrl *string    `json:"pricing_page_url,omitempty"`
	Labels         []string   `json:"labels" binding:"required"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

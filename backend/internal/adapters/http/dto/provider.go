package dto

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

// ProviderModel represents a service provider offering subscription plans
// @Description Provider object containing information about a subscription service provider and their available plans
type ProviderModel struct {
	// @Description Unique identifier for the provider (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the service provider
	Name string  `json:"name" binding:"required" example:"Netflix" minLength:"1" maxLength:"255"`
	Key  *string `json:"key,omitempty"`
	// @Description Optional detailed description of the provider and their services
	Description *string `json:"description,omitempty" example:"Streaming service offering movies and TV shows"`
	// @Description Optional URL to the provider's icon or logo image
	IconUrl *string `json:"icon_url,omitempty" example:"https://example.com/netflix-icon.png"`
	// @Description Optional URL to the provider's main website
	Url *string `json:"url,omitempty" example:"https://netflix.com"`
	// @Description Optional URL to the provider's pricing information page
	PricingPageUrl *string `json:"pricing_page_url,omitempty" example:"https://netflix.com/pricing"`
	// @Description List of label IDs associated with this provider for categorization
	Labels []string `json:"labels" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001,123e4567-e89b-12d3-a456-426614174002"`
	// @Description Ownership information specifying whether this provider belongs to a user or family
	Owner OwnerModel `json:"owner" binding:"required"`
	// @Description ISO 8601 timestamp when the provider was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp when the provider was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func NewProviderModel(source provider.Provider) ProviderModel {
	return ProviderModel{
		Id:             source.Id().String(),
		Name:           source.Name(),
		Key:            source.Key(),
		Description:    source.Description(),
		IconUrl:        source.IconUrl(),
		Url:            source.Url(),
		PricingPageUrl: source.PricingPageUrl(),
		Labels:         herd.Select(source.Labels().Values(), func(id types.LabelID) string { return id.String() }),
		Owner:          NewOwnerModel(source.Owner()),
		CreatedAt:      source.CreatedAt(),
		UpdatedAt:      source.UpdatedAt(),
		Etag:           source.ETag(),
	}
}

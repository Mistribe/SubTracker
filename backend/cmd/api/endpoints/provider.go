package endpoints

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type ProviderEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func (s ProviderEndpointGroup) Prefix() string {
	return "/providers"
}

func (s ProviderEndpointGroup) Routes() []ginfx.Route {
	return s.routes
}

func (s ProviderEndpointGroup) Middlewares() []gin.HandlerFunc {
	return s.middlewares
}

func NewProviderEndpointGroup(
	getEndpoint *ProviderGetEndpoint,
	getAllEndpoint *ProviderGetAllEndpoint,
	createEndpoint *ProviderCreateEndpoint,
	createPlanEndpoint *ProviderPlanCreateEndpoint,
	createPriceEndpoint *ProviderPriceCreateEndpoint,
	updateEndpoint *ProviderUpdateEndpoint,
	updatePlanEndpoint *ProviderPlanUpdateEndpoint,
	updatePriceEndpoint *ProviderPriceUpdateEndpoint,
	deleteEndpoint *ProviderDeleteEndpoint,
	deletePlanEndpoint *ProviderPlanDeleteEndpoint,
	deletePriceEndpoint *ProviderPriceDeleteEndpoint,
	patchProviderEndpoint *ProviderPatchEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *ProviderEndpointGroup {
	return &ProviderEndpointGroup{
		routes: []ginfx.Route{
			getEndpoint,
			getAllEndpoint,
			createEndpoint,
			createPlanEndpoint,
			createPriceEndpoint,
			updateEndpoint,
			updatePlanEndpoint,
			updatePriceEndpoint,
			deleteEndpoint,
			deletePlanEndpoint,
			deletePriceEndpoint,
			patchProviderEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

// PriceModel represents a pricing tier with currency and validity period
// @Description Price object defining the cost of a subscription plan with currency and time validity
type PriceModel struct {
	// @Description Unique identifier for the price (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description ISO 4217 currency code for the price
	Currency string `json:"currency" binding:"required" example:"USD" minLength:"3" maxLength:"3"`
	// @Description ISO 8601 timestamp when this price becomes effective
	StartDate time.Time `json:"start_date" binding:"required" format:"date-time" example:"2023-01-01T00:00:00Z"`
	// @Description ISO 8601 timestamp when this price expires (null means indefinite)
	EndDate *time.Time `json:"end_date,omitempty" format:"date-time" example:"2023-12-31T23:59:59Z"`
	// @Description Price amount in the specified currency (supports decimal values)
	Amount float64 `json:"amount" binding:"required" example:"9.99" minimum:"0"`
	// @Description ISO 8601 timestamp when the price was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp when the price was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

// PlanModel represents a subscription plan offered by a provider
// @Description Plan object defining a specific subscription tier with associated pricing options
type PlanModel struct {
	// @Description Unique identifier for the plan (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the subscription plan
	Name string `json:"name" binding:"required" example:"Premium Plan" minLength:"1" maxLength:"255"`
	// @Description Optional detailed description of the plan features and benefits
	Description *string `json:"description,omitempty" example:"Access to all premium features with unlimited usage"`
	// @Description List of pricing options available for this plan (different currencies, time periods, etc.)
	Prices []PriceModel `json:"prices" binding:"required"`
	// @Description ISO 8601 timestamp when the plan was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp when the plan was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

// ProviderModel represents a service provider offering subscription plans
// @Description Provider object containing information about a subscription service provider and their available plans
type ProviderModel struct {
	// @Description Unique identifier for the provider (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the service provider
	Name string `json:"name" binding:"required" example:"Netflix" minLength:"1" maxLength:"255"`
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
	// @Description List of subscription plans offered by this provider
	Plans []PlanModel `json:"plans" binding:"required"`
	// @Description Ownership information specifying whether this provider belongs to a user or family
	Owner OwnerModel `json:"owner" binding:"required"`
	// @Description ISO 8601 timestamp when the provider was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp when the provider was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func newProviderModel(source provider.Provider) ProviderModel {
	return ProviderModel{
		Id:             source.Id().String(),
		Name:           source.Name(),
		Description:    source.Description(),
		IconUrl:        source.IconUrl(),
		Url:            source.Url(),
		PricingPageUrl: source.PricingPageUrl(),
		Labels:         slicesx.Select(source.Labels().Values(), func(id uuid.UUID) string { return id.String() }),
		Plans:          slicesx.Select(source.Plans().Values(), newPlanModel),
		Owner:          newOwnerModel(source.Owner()),
		CreatedAt:      source.CreatedAt(),
		UpdatedAt:      source.UpdatedAt(),
		Etag:           source.ETag(),
	}
}

func newPlanModel(source provider.Plan) PlanModel {
	return PlanModel{
		Id:          source.Id().String(),
		Name:        source.Name(),
		Description: source.Description(),
		Prices:      slicesx.Select(source.Prices().Values(), newPriceModel),
		CreatedAt:   source.CreatedAt(),
		UpdatedAt:   source.UpdatedAt(),
		Etag:        source.ETag(),
	}
}

func newPriceModel(source provider.Price) PriceModel {
	return PriceModel{
		Id:        source.Id().String(),
		Currency:  source.Currency().String(),
		StartDate: source.StartDate(),
		EndDate:   source.EndDate(),
		Amount:    source.Amount(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

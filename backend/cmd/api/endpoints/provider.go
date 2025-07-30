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

type PriceModel struct {
	Id        string     `json:"id" binding:"required"`
	Currency  string     `json:"currency" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Amount    float64    `json:"amount" binding:"required"`
	CreatedAt time.Time  `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt time.Time  `json:"updated_at" binding:"required" format:"date-time"`
	Etag      string     `json:"etag" binding:"required"`
}

type PlanModel struct {
	Id          string       `json:"id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description *string      `json:"description,omitempty"`
	Prices      []PriceModel `json:"prices" binding:"required"`
	CreatedAt   time.Time    `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt   time.Time    `json:"updated_at" binding:"required" format:"date-time"`
	Etag        string       `json:"etag" binding:"required"`
}

type ProviderModel struct {
	Id             string      `json:"id" binding:"required"`
	Name           string      `json:"name" binding:"required"`
	Description    *string     `json:"description,omitempty"`
	IconUrl        *string     `json:"icon_url,omitempty"`
	Url            *string     `json:"url,omitempty"`
	PricingPageUrl *string     `json:"pricing_page_url,omitempty"`
	Labels         []string    `json:"labels" binding:"required"`
	Plans          []PlanModel `json:"plans" binding:"required"`
	Owner          OwnerModel  `json:"owner" binding:"required"`
	CreatedAt      time.Time   `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt      time.Time   `json:"updated_at" binding:"required" format:"date-time"`
	Etag           string      `json:"etag" binding:"required"`
}

func newProviderModel(source provider.Provider) ProviderModel {
	return ProviderModel{
		Id:             source.Id().String(),
		Name:           source.Name(),
		Description:    source.Description(),
		IconUrl:        source.IconUrl(),
		Url:            source.Url(),
		PricingPageUrl: source.PricingPageUrl(),
		Labels:         slicesx.Map(source.Labels().Values(), func(id uuid.UUID) string { return id.String() }),
		Plans:          slicesx.Map(source.Plans().Values(), newPlanModel),
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
		Prices:      slicesx.Map(source.Prices().Values(), newPriceModel),
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

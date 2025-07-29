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
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
			getEndpoint,
			getAllEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

type priceModel struct {
	Id        string     `json:"id" binding:"required"`
	Currency  string     `json:"currency" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Amount    float64    `json:"amount" binding:"required"`
	CreatedAt string     `json:"created_at" binding:"required"`
	UpdatedAt string     `json:"updated_at" binding:"required"`
	Etag      string     `json:"etag" binding:"required"`
}
type planModel struct {
	Id          string       `json:"id" binding:"required"`
	Name        *string      `json:"name,omitempty"`
	Description *string      `json:"description,omitempty"`
	Prices      []priceModel `json:"prices" binding:"required"`
	CreatedAt   string       `json:"created_at" binding:"required"`
	UpdatedAt   string       `json:"updated_at" binding:"required"`
	Etag        string       `json:"etag" binding:"required"`
}
type providerModel struct {
	Id             string      `json:"id" binding:"required"`
	Name           string      `json:"name" binding:"required"`
	Description    *string     `json:"description,omitempty"`
	IconUrl        *string     `json:"icon_url,omitempty"`
	Url            *string     `json:"url,omitempty"`
	PricingPageUrl *string     `json:"pricing_page_url,omitempty"`
	Labels         []string    `json:"labels" binding:"required"`
	Plans          []planModel `json:"plans" binding:"required"`
	Owner          *ownerModel `json:"owner,omitempty"`
	CreatedAt      string      `json:"created_at" binding:"required"`
	UpdatedAt      string      `json:"updated_at" binding:"required"`
	Etag           string      `json:"etag" binding:"required"`
}

func newProviderModel(source provider.Provider) providerModel {

	model := providerModel{
		Id:             source.Id().String(),
		Name:           source.Name(),
		Description:    source.Description(),
		IconUrl:        source.IconUrl(),
		Url:            source.Url(),
		PricingPageUrl: source.PricingPageUrl(),
		Labels:         slicesx.Map(source.Labels().Values(), func(id uuid.UUID) string { return id.String() }),
		Plans:          slicesx.Map(source.Plans().Values(), newPlanModel),
		CreatedAt:      source.CreatedAt().String(),
		UpdatedAt:      source.UpdatedAt().String(),
		Etag:           source.ETag(),
	}

	if source.Owner() != nil {
		owner := newOwnerModel(source.Owner())
		model.Owner = &owner
	}

	return model
}

func newPlanModel(source provider.Plan) planModel {
	return planModel{
		Id:          source.Id().String(),
		Name:        source.Name(),
		Description: source.Description(),
		Prices:      slicesx.Map(source.Prices().Values(), newPriceModel),
		CreatedAt:   source.CreatedAt().String(),
		UpdatedAt:   source.UpdatedAt().String(),
		Etag:        source.ETag(),
	}
}

func newPriceModel(source provider.Price) priceModel {
	return priceModel{
		Id:        source.Id().String(),
		Currency:  source.Currency().String(),
		StartDate: source.StartDate(),
		EndDate:   source.EndDate(),
		Amount:    source.Amount(),
		CreatedAt: source.CreatedAt().String(),
		UpdatedAt: source.UpdatedAt().String(),
		Etag:      source.ETag(),
	}
}

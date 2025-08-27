package label

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
	"github.com/mistribe/subtracker/internal/domain/label"
)

type EndpointGroup struct {
	routes      []ginfx.Endpoint
	middlewares []gin.HandlerFunc
}

func NewEndpointGroup(
	createEndpoint *LabelCreateEndpoint,
	updateEndpoint *LabelUpdateEndpoint,
	deleteEndpoint *LabelDeleteEndpoint,
	getEndpoint *LabelGetEndpoint,
	getAllEndpoint *LabelGetAllEndpoint,
	defaultEndpoint *DefaultLabelEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			getEndpoint,
			getAllEndpoint,
			defaultEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g EndpointGroup) Prefix() string {
	return "/labels"
}

func (g EndpointGroup) Routes() []ginfx.Endpoint {
	return g.routes
}

func (g EndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}

// Label represents a categorization tag for organizing subscriptions and expenses
// @Description Label object used for categorizing and organizing subscriptions with customizable colors
type labelModel struct {
	// @Description Unique identifier for the label (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the label
	Name string  `json:"name" binding:"required" example:"Entertainment" minLength:"1" maxLength:"100"`
	Key  *string `json:"key,omitempty"`
	// @Description Hexadecimal color code for visual representation of the label
	Color string `json:"color" binding:"required" example:"#FF5733" pattern:"^#[0-9A-Fa-f]{6}$"`
	// @Description Ownership information specifying whether this label belongs to a user or family
	Owner dto.OwnerModel `json:"owner" binding:"required"`
	// @Description ISO 8601 timestamp indicating when the label was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp indicating when the label was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func newLabelModel(source label.Label) labelModel {
	return labelModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		Key:       source.Key(),
		Color:     source.Color(),
		Owner:     dto.NewOwnerModel(source.Owner()),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

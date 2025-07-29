package endpoints

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func NewLabelEndpointGroup(
	createEndpoint *LabelCreateEndpoint,
	updateEndpoint *LabelUpdateEndpoint,
	deleteEndpoint *LabelDeleteEndpoint,
	getEndpoint *LabelGetEndpoint,
	getAllEndpoint *LabelGetAllEndpoint,
	defaultEndpoint *DefaultLabelEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *LabelEndpointGroup {
	return &LabelEndpointGroup{
		routes: []ginfx.Route{
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

func (g LabelEndpointGroup) Prefix() string {
	return "/labels"
}

func (g LabelEndpointGroup) Routes() []ginfx.Route {
	return g.routes
}

func (g LabelEndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}

type labelModel struct {
	Id        string      `json:"id" binding:"required"`
	Name      string      `json:"name" binding:"required"`
	IsDefault bool        `json:"is_default" binding:"required"`
	Color     string      `json:"color" binding:"required"`
	Owner     *ownerModel `json:"owner,omitempty"`
	CreatedAt time.Time   `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt time.Time   `json:"updated_at" binding:"required" format:"date-time"`
	Etag      string      `json:"etag" binding:"required"`
}

func newLabelModel(source label.Label) labelModel {
	model := labelModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsDefault: source.IsDefault(),
		Color:     source.Color(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}

	if source.Owner() != nil {
		owner := newOwnerModel(source.Owner())
		model.Owner = &owner
	}

	return model
}

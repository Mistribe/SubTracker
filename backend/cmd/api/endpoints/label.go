package endpoints

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelEndpointGroup struct {
	routes []ginfx.Route
}

func NewLabelEndpointGroup(
	createEndpoint *LabelCreateEndpoint,
	updateEndpoint *LabelUpdateEndpoint,
	deleteEndpoint *LabelDeleteEndpoint,
	getEndpoint *LabelGetEndpoint,
	getAllEndpoint *LabelGetAllEndpoint) *LabelEndpointGroup {
	return &LabelEndpointGroup{
		routes: []ginfx.Route{
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			getEndpoint,
			getAllEndpoint,
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
	return nil
}

type labelModel struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"is_default"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newLabelModel(source label.Label) labelModel {
	return labelModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsDefault: source.IsDefault(),
		Color:     source.Color(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
	}
}

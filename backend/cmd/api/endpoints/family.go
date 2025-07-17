package endpoints

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func NewFamilyEndpointGroup(
	createEndpoint *FamilyMemberCreateEndpoint,
	updateEndpoint *FamilyMemberUpdateEndpoint,
	deleteEndpoint *FamilyMemberDeleteEndpoint,
	getEndpoint *FamilyMemberGetEndpoint,
	getAllEndpoint *FamilyMemberGetAllEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *FamilyEndpointGroup {
	return &FamilyEndpointGroup{
		routes: []ginfx.Route{
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			getEndpoint,
			getAllEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g FamilyEndpointGroup) Prefix() string {
	return "/families"
}

func (g FamilyEndpointGroup) Routes() []ginfx.Route {
	return g.routes
}

func (g FamilyEndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}

type familyMemberModel struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	IsKid     bool      `json:"is_kid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newFamilyMemberModel(source family.Member) familyMemberModel {
	return familyMemberModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsKid:     source.IsKid(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
	}
}

package endpoints

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyEndpointGroup struct {
	routes []ginfx.Route
}

func NewFamilyEndpointGroup(
	createEndpoint *FamilyMemberCreateEndpoint,
	updateEndpoint *FamilyMemberUpdateEndpoint,
	deleteEndpoint *FamilyMemberDeleteEndpoint,
	getEndpoint *FamilyMemberGetEndpoint,
	getAllEndpoint *FamilyMemberGetAllEndpoint) *FamilyEndpointGroup {
	return &FamilyEndpointGroup{
		routes: []ginfx.Route{
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			getEndpoint,
			getAllEndpoint,
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
	return nil
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

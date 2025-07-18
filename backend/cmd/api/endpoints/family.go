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
	familyCreateEndpoint *FamilyCreateEndpoint,
	familyUpdateEndpoint *FamilyUpdateEndpoint,
	familyMemberCreateEndpoint *FamilyMemberCreateEndpoint,
	familyMemberUpdateEndpoint *FamilyMemberUpdateEndpoint,
	familyMemberDeleteEndpoint *FamilyMemberDeleteEndpoint,
	familyGetEndpoint *FamilyGetEndpoint,
	familyGetAllEndpoint *FamilyGetAllEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *FamilyEndpointGroup {
	return &FamilyEndpointGroup{
		routes: []ginfx.Route{
			familyCreateEndpoint,
			familyUpdateEndpoint,
			familyMemberCreateEndpoint,
			familyMemberUpdateEndpoint,
			familyMemberDeleteEndpoint,
			familyGetEndpoint,
			familyGetAllEndpoint,
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

type familyModel struct {
	Id               string              `json:"id"`
	Name             string              `json:"name"`
	IsOwner          bool                `json:"is_owner"`
	Members          []familyMemberModel `json:"members"`
	HaveJointAccount bool                `json:"have_joint_account"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

func newFamilyModel(userId string, source family.Family) familyModel {
	members := make([]familyMemberModel, 0, len(source.Members()))
	for _, member := range source.Members() {
		members = append(members, newFamilyMemberModel(member))
	}

	return familyModel{
		Id:               source.Id().String(),
		Name:             source.Name(),
		IsOwner:          source.OwnerId() == userId,
		Members:          members,
		HaveJointAccount: source.HaveJointAccount(),
		CreatedAt:        source.CreatedAt(),
		UpdatedAt:        source.UpdatedAt(),
	}
}

type familyMemberModel struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	IsKid     bool      `json:"is_kid"`
	FamilyId  string    `json:"family_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newFamilyMemberModel(source family.Member) familyMemberModel {
	return familyMemberModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsKid:     source.IsKid(),
		FamilyId:  source.FamilyId().String(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
	}
}

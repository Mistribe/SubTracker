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
	familyPatchEndpoint *FamilyPatchEndpoint,
	familyDeleteEndpoint *FamilyDeleteEndpoint,
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
			familyPatchEndpoint,
			familyDeleteEndpoint,
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
	Id               string              `json:"id" binding:"required"`
	Name             string              `json:"name" binding:"required"`
	IsOwner          bool                `json:"is_owner" binding:"required"`
	Members          []familyMemberModel `json:"members" binding:"required"`
	HaveJointAccount bool                `json:"have_joint_account" binding:"required"`
	CreatedAt        time.Time           `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt        time.Time           `json:"updated_at" binding:"required" format:"date-time"`
	Etag             string              `json:"etag" binding:"required"`
}

func newFamilyModel(userId string, source family.Family) familyModel {
	var members []familyMemberModel
	for member := range source.Members().It() {
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
		Etag:             source.ETag(),
	}
}

type familyMemberModel struct {
	Id        string    `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IsKid     bool      `json:"is_kid" binding:"required"`
	FamilyId  string    `json:"family_id" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time"`
	Etag      string    `json:"etag" binding:"required"`
}

func newFamilyMemberModel(source family.member) familyMemberModel {
	return familyMemberModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsKid:     source.IsKid(),
		FamilyId:  source.FamilyId().String(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

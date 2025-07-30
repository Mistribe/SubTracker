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

// Family represents a family unit with its members and settings
// @Description Family object containing family information and members
type familyModel struct {
	// @Description Unique identifier for the family (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the family
	Name string `json:"name" binding:"required" example:"Smith Family" minLength:"1" maxLength:"255"`
	// @Description Indicates whether the current authenticated user is the owner of this family
	IsOwner bool `json:"is_owner" binding:"required" example:"true"`
	// @Description Complete list of all members belonging to this family
	Members []familyMemberModel `json:"members" binding:"required"`
	// @Description ISO 8601 timestamp indicating when the family was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp indicating when the family information was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func newFamilyModel(userId string, source family.Family) familyModel {
	var members []familyMemberModel
	for member := range source.Members().It() {
		members = append(members, newFamilyMemberModel(member))
	}

	return familyModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsOwner:   source.OwnerId() == userId,
		Members:   members,
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

// FamilyMember represents a member of a family
// @Description Family member object containing member information
type familyMemberModel struct {
	// @Description Unique identifier for the family member
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
	// @Description Name of the family member
	Name string `json:"name" binding:"required" example:"John Smith"`
	// @Description Whether this member is a child (affects permissions and features)
	IsKid bool `json:"is_kid" binding:"required" example:"false"`
	// @Description ID of the family this member belongs to
	FamilyId string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Timestamp when the member was created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time"`
	// @Description Timestamp when the member was last updated
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time"`
	// @Description Entity tag for optimistic concurrency control
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func newFamilyMemberModel(source family.Member) familyMemberModel {
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

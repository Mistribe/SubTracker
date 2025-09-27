package dto

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/family"
)

// Family represents a family unit with its members and settings
// @Description Family object containing family information and members
type FamilyModel struct {
	// @Description Unique identifier for the family (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the family
	Name string `json:"name" binding:"required" example:"Smith Family" minLength:"1" maxLength:"255"`
	// @Description Indicates whether the current authenticated user is the owner of this family
	IsOwner bool `json:"is_owner" binding:"required" example:"true"`
	// @Description Complete list of all members belonging to this family
	Members []FamilyMemberModel `json:"members" binding:"required"`
	// @Description ISO 8601 timestamp indicating when the family was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp indicating when the family information was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func NewFamilyModel(userId account.UserID, source family.Family) FamilyModel {
	var members []FamilyMemberModel
	for member := range source.Members().It() {
		members = append(members, NewFamilyMemberModel(userId, member))
	}

	return FamilyModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		IsOwner:   source.Owner().UserId() == userId,
		Members:   members,
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

// FamilyMember represents a member of a family
// @Description Family member object containing member information
type FamilyMemberModel struct {
	// @Description Unique identifier for the family member
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
	// @Description Name of the family member
	Name string `json:"name" binding:"required" example:"John Smith"`
	// @Description Whether this member is a child (affects permissions and features)
	Type string `json:"type" binding:"required" enums:"owner,adult,kid"`
	// @Description LabelID of the family this member belongs to
	FamilyId string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Indicates whether this member is the current authenticated user
	IsYou bool `json:"is_you" binding:"required" example:"false"`
	// @Description Indicates whether this member has an account with the service provider
	HasAccount bool `json:"has_account" binding:"required" example:"true"`
	// @Description Timestamp when the member was created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time"`
	// @Description Timestamp when the member was last updated
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time"`
	// @Description Entity tag for optimistic concurrency control
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func NewFamilyMemberModel(currentUserId string, source family.Member) FamilyMemberModel {
	model := FamilyMemberModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		Type:      source.Type().String(),
		FamilyId:  source.FamilyId().String(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}

	if source.UserId() != nil {
		model.IsYou = *source.UserId() == currentUserId
		model.HasAccount = true
	}

	return model
}

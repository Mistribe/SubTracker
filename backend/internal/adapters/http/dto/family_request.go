package dto

import (
	"time"
)

// FamilyAcceptInvitationRequest represents the request body for accepting a family invitation
type FamilyAcceptInvitationRequest struct {
	// Code received in the invitation
	InvitationCode string `json:"invitation_code" binding:"required" example:"123456"`
	// ID of the family member accepting the invitation
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type CreateFamilyRequest struct {
	Id          *string    `json:"id,omitempty"`
	Name        string     `json:"name" binding:"required"`
	CreatorName string     `json:"creator_name" binding:"required"`
	CreatedAt   *time.Time `json:"created_at,omitempty" format:"date-time"`
}

type FamilyDeclineInvitationRequest struct {
	// Code received in the invitation
	InvitationCode string `json:"invitation_code" binding:"required" example:"123456"`
	// ID of the family member accepting the invitation
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// FamilyInviteRequest represents the request body for inviting a family member
type FamilyInviteRequest struct {
	// Email of the invited member
	Email *string `json:"email,omitempty"`
	// ID of the family member to be invited
	FamilyMemberId string `json:"family_member_id" binding:"required"`
	// Name of the invited member
	Name *string `json:"name,omitempty"`
	// Type of the member (adult or kid)
	Type *string `json:"type,omitempty" enums:"adult,kid"`
}

type CreateFamilyMemberRequest struct {
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name" binding:"required"`
	Type      string     `json:"type" binding:"required" enums:"owner,adult,kid"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

type UpdateFamilyMemberRequest struct {
	Name      string     `json:"name" binding:"required"`
	Type      string     `json:"type" binding:"required" enums:"owner,adult,kid"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

type UpdateFamilyRequest struct {
	Name      string     `json:"name" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

package dto

type FamilyInviteResponse struct {
	Code           string `json:"code" binding:"required" example:"123456"`
	FamilyId       string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
}

type FamilySeeInvitationResponse struct {
	// Family details
	Family FamilyModel `json:"family"`
	// Role of the invited member
	InvitedInasmuchAs string `json:"invited_inasmuch_as" example:"OWNER"`
}

type UserFamilyResponse struct {
	Family FamilyModel `json:"family"`
	Limits []Limit     `json:"limits,omitempty"`
}

package family

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrFamilyMemberNotFound      = ex.NewNotFound("family member not found")
	ErrFamilyNotFound            = ex.NewNotFound("family not found")
	ErrFamilyAlreadyExists       = ex.NewAlreadyExists("family already exists")
	ErrOnlyOwnerCanEditFamily    = ex.NewUnauthorized("only the owner can edit the family")
	ErrDuplicateMember           = ex.NewAlreadyExists("member with the same LabelID or email already exists")
	ErrCannotInviteUser          = ex.NewUnauthorized("cannot invite user to family")
	ErrBadInvitationCode         = ex.NewInvalidValue("bad invitation code")
	ErrFamilyMembersLimitReached = ex.NewInvalidValue("family members limit reached")
)

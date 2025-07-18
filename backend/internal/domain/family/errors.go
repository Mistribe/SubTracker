package family

import "errors"

var (
	ErrMemberAlreadyExists    = errors.New("family member already exists")
	ErrFamilyMemberNotFound   = errors.New("family member not found")
	ErrFamilyNotFound         = errors.New("family not found")
	ErrNameRequired           = errors.New("name is required")
	ErrNameTooShort           = errors.New("name must be at least 3 characters")
	ErrNameTooLong            = errors.New("name must be less than 100 characters")
	ErrFamilyAlreadyExists    = errors.New("family already exists")
	ErrOnlyOwnerCanEditFamily = errors.New("only the owner can edit the family")
)

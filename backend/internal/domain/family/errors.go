package family

import "errors"

var (
	ErrMemberAlreadyExists  = errors.New("family member already exists")
	ErrFamilyMemberNotFound = errors.New("family member not found")
)

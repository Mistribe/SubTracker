package types

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/entity"
)

var (
	SystemOwner         = NewSystemOwner()
	ErrUnknownOwnerType = errors.New("unknown owner type")
)

type OwnerType string
type OwnerTypes []OwnerType

const (
	UnknownOwnerType  OwnerType = "unknown"
	SystemOwnerType   OwnerType = "system"
	PersonalOwnerType OwnerType = "personal"
	FamilyOwnerType   OwnerType = "family"
)

func (o OwnerType) String() string {
	return string(o)
}

func (o OwnerTypes) Contains(t OwnerType) bool {
	return slices.Contains(o, t)
}

func (o OwnerTypes) Strings() []string {
	var result []string
	for _, ot := range o {
		result = append(result, ot.String())
	}
	return result
}

func ParseOwnerType(input string) (OwnerType, error) {
	switch input {
	case string(PersonalOwnerType):
		return PersonalOwnerType, nil
	case string(FamilyOwnerType):
		return FamilyOwnerType, nil
	case string(SystemOwnerType):
		return SystemOwnerType, nil
	default:
		return UnknownOwnerType, ErrUnknownOwnerType
	}
}

func MustParseOwnerType(input string) OwnerType {
	t, err := ParseOwnerType(input)
	if err != nil {
		panic(err)
	}
	return t
}

type Owner interface {
	entity.ETagEntity

	FamilyId() FamilyID
	UserId() UserID
	Type() OwnerType
	Equal(other Owner) bool
}

func NewOwner(ownerType OwnerType, familyId *uuid.UUID, userId *string) Owner {
	switch ownerType {
	case PersonalOwnerType:
		if userId == nil {
			panic("missing userProfile id for a personal owner type")
		}
		return NewPersonalOwner(UserID(*userId))
	case FamilyOwnerType:
		if familyId == nil {
			panic("missing family id for a family owner type")
		}
		return NewFamilyOwner(*familyId)
	default:
		return NewSystemOwner()
	}
}

type familyOwner struct {
	familyId FamilyID
}

func NewFamilyOwner(familyId FamilyID) Owner {
	return &familyOwner{
		familyId: familyId,
	}
}

func (o familyOwner) FamilyId() FamilyID {
	return o.familyId
}

func (o familyOwner) UserId() UserID {
	panic("family owner cannot have owner id")
}

func (o familyOwner) Type() OwnerType {
	return FamilyOwnerType
}

func (o familyOwner) Equal(other Owner) bool {
	if other == nil {
		return false
	}

	return o.ETag() == other.ETag()
}

func (o familyOwner) ETagFields() []interface{} {
	return []interface{}{
		o.familyId.String(),
		FamilyOwnerType.String(),
	}

}
func (o familyOwner) ETag() string {
	return entity.CalculateETag(o)
}

type personalOwner struct {
	ownerId UserID
}

func NewPersonalOwner(ownerId UserID) Owner {
	return &personalOwner{
		ownerId: ownerId,
	}
}

func (o personalOwner) FamilyId() FamilyID {
	panic("personal owner cannot have family id")
}

func (o personalOwner) UserId() UserID {
	return o.ownerId
}

func (o personalOwner) Type() OwnerType {
	return PersonalOwnerType
}

func (o personalOwner) Equal(other Owner) bool {
	if other == nil {
		return false
	}

	return o.ETag() == other.ETag()
}

func (o personalOwner) ETagFields() []interface{} {
	return []interface{}{
		o.ownerId,
		PersonalOwnerType.String(),
	}

}
func (o personalOwner) ETag() string {
	return entity.CalculateETag(o)
}

type systemOwner struct {
}

func NewSystemOwner() Owner {
	return &systemOwner{}
}

func (o systemOwner) FamilyId() FamilyID {
	panic("system owner cannot have family id")
}

func (o systemOwner) UserId() UserID {
	panic("system owner cannot have userProfile id")
}

func (o systemOwner) Type() OwnerType {
	return SystemOwnerType
}

func (o systemOwner) Equal(other Owner) bool {
	if other == nil {
		return false
	}

	return o.ETag() == other.ETag()
}

func (o systemOwner) ETagFields() []interface{} {
	return []interface{}{
		PersonalOwnerType.String(),
	}

}
func (o systemOwner) ETag() string {
	return entity.CalculateETag(o)
}

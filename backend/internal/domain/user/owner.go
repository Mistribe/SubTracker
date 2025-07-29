package user

import (
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type OwnerType string

const (
	UnknownOwner  OwnerType = "unknown"
	SystemOwner   OwnerType = "system"
	PersonalOwner OwnerType = "personal"
	FamilyOwner   OwnerType = "family"
)

func (o OwnerType) String() string {
	return string(o)
}

func ParseOwnerType(input string) (OwnerType, error) {
	switch input {
	case string(PersonalOwner):
		return PersonalOwner, nil
	case string(FamilyOwner):
		return FamilyOwner, nil
	case string(SystemOwner):
		return SystemOwner, nil
	default:
		return UnknownOwner, ErrUnknownOwnerType
	}
}

type Owner interface {
	entity.ETagEntity

	FamilyId() uuid.UUID
	UserId() string
	Type() OwnerType
	Equal(other Owner) bool
}

func NewOwner(ownerType OwnerType, familyId *uuid.UUID, userId *string) Owner {
	switch ownerType {
	case PersonalOwner:
		if userId == nil {
			panic("missing user id for a personal owner type")
		}
		return NewPersonalOwner(*userId)
	case FamilyOwner:
		if familyId == nil {
			panic("missing family id for a family owner type")
		}
		return NewFamilyOwner(*familyId)
	default:
		return NewSystemOwner()
	}
}

type familyOwner struct {
	familyId uuid.UUID
}

func NewFamilyOwner(familyId uuid.UUID) Owner {
	return &familyOwner{
		familyId: familyId,
	}
}

func (o familyOwner) FamilyId() uuid.UUID {
	return o.familyId
}

func (o familyOwner) UserId() string {
	panic("family owner cannot have owner id")
}

func (o familyOwner) Type() OwnerType {
	return FamilyOwner
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
		FamilyOwner.String(),
	}

}
func (o familyOwner) ETag() string {
	return entity.CalculateETag(o)
}

type personalOwner struct {
	ownerId string
}

func NewPersonalOwner(ownerId string) Owner {
	return &personalOwner{
		ownerId: ownerId,
	}
}

func (o personalOwner) FamilyId() uuid.UUID {
	panic("personal owner cannot have family id")
}

func (o personalOwner) UserId() string {
	return o.ownerId
}

func (o personalOwner) Type() OwnerType {
	return PersonalOwner
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
		PersonalOwner.String(),
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

func (o systemOwner) FamilyId() uuid.UUID {
	panic("system owner cannot have family id")
}

func (o systemOwner) UserId() string {
	panic("system owner cannot have user id")
}

func (o systemOwner) Type() OwnerType {
	return SystemOwner
}

func (o systemOwner) Equal(other Owner) bool {
	if other == nil {
		return false
	}

	return o.ETag() == other.ETag()
}

func (o systemOwner) ETagFields() []interface{} {
	return []interface{}{
		PersonalOwner.String(),
	}

}
func (o systemOwner) ETag() string {
	return entity.CalculateETag(o)
}

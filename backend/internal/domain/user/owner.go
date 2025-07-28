package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type OwnerType string

const (
	UnknownOwner  OwnerType = "unknown"
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
	default:
		return UnknownOwner, ErrUnknownOwnerType
	}
}

type Owner interface {
	entity.Entity

	FamilyId() uuid.UUID
	OwnerId() string
	Type() OwnerType
	Equal(other Owner) bool
}

type familyOwner struct {
	*entity.Base

	familyId uuid.UUID
}

func NewFamilyOwner(
	id uuid.UUID,
	familyId uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time) Owner {
	return &familyOwner{
		Base:     entity.NewBase(id, createdAt, updatedAt, true, false),
		familyId: familyId,
	}
}

func (o familyOwner) FamilyId() uuid.UUID {
	return o.familyId
}

func (o familyOwner) OwnerId() string {
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
	return entity.CalculateETag(o, o.Base)
}

type personalOwner struct {
	*entity.Base

	ownerId string
}

func NewPersonalOwner(
	id uuid.UUID,
	ownerId string,
	createdAt time.Time,
	updatedAt time.Time) Owner {
	return &personalOwner{
		Base:    entity.NewBase(id, createdAt, updatedAt, true, false),
		ownerId: ownerId,
	}
}

func (o personalOwner) FamilyId() uuid.UUID {
	panic("personal owner cannot have family id")
}

func (o personalOwner) OwnerId() string {
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
	return entity.CalculateETag(o, o.Base)
}

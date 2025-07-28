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
	entity.Entity

	familyId uuid.UUID
}

func NewFamilyOwner(
	id uuid.UUID,
	familyId uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time) Owner {
	return &familyOwner{
		Entity:   entity.NewBase(id, createdAt, updatedAt, true, false),
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
	if other.Type() != FamilyOwner {
		return false
	}
	return o.Id() == other.Id() &&
		o.familyId == other.FamilyId() &&
		o.CreatedAt() == other.CreatedAt() &&
		o.UpdatedAt() == other.UpdatedAt()
}

type personalOwner struct {
	entity.Entity

	ownerId string
}

func NewPersonalOwner(
	id uuid.UUID,
	ownerId string,
	createdAt time.Time,
	updatedAt time.Time) Owner {
	return &personalOwner{
		Entity:  entity.NewBase(id, createdAt, updatedAt, true, false),
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
	if other.Type() != PersonalOwner {
		return false
	}
	return o.Id() == other.Id() &&
		o.ownerId == other.OwnerId() &&
		o.CreatedAt() == other.CreatedAt() &&
		o.UpdatedAt() == other.UpdatedAt()
}

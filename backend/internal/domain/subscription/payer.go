package subscription

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type PayerType string

const (
	UnknownPayer      PayerType = "unknown"
	FamilyMemberPayer PayerType = "family_member"
	FamilyPayer       PayerType = "family"
)

func ParsePayerType(input string) (payerType PayerType, err error) {
	switch input {
	case string(FamilyMemberPayer):
		return FamilyMemberPayer, nil
	case string(FamilyPayer):
		return FamilyPayer, nil
	default:
		return UnknownPayer, ErrUnknownPayerType
	}
}

type Payer interface {
	entity.Entity
	entity.ETagEntity

	Type() PayerType
	FamilyId() uuid.UUID
	MemberId() uuid.UUID
}

type familyPayer struct {
	*entity.Base

	familyId uuid.UUID
}

func NewFamilyPayer(id uuid.UUID,
	familyId uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time) Payer {
	return familyPayer{
		Base:     entity.NewBase(id, createdAt, updatedAt, true, false),
		familyId: familyId,
	}
}

func (f familyPayer) Type() PayerType {
	return FamilyPayer
}

func (f familyPayer) FamilyId() uuid.UUID {
	return f.familyId
}

func (f familyPayer) MemberId() uuid.UUID {
	panic("family payer cannot have member id")
}

type familyMemberPayer struct {
	*entity.Base

	familyId uuid.UUID
	memberId uuid.UUID
}

func NewFamilyMemberPayer(id uuid.UUID,
	familyId uuid.UUID,
	memberId uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time) Payer {
	return familyMemberPayer{
		Base:     entity.NewBase(id, createdAt, updatedAt, true, false),
		familyId: familyId,
		memberId: memberId,
	}
}

func (f familyMemberPayer) FamilyId() uuid.UUID {
	return f.familyId
}

func (f familyMemberPayer) MemberId() uuid.UUID {
	return f.memberId
}

func (f familyMemberPayer) Type() PayerType {
	return FamilyMemberPayer
}

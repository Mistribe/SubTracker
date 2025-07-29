package subscription

import (
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type PayerType string

const (
	UnknownPayer      PayerType = "unknown"
	FamilyMemberPayer PayerType = "family_member"
	FamilyPayer       PayerType = "family"
)

func (p PayerType) String() string {
	return string(p)
}

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
	entity.ETagEntity

	Type() PayerType
	FamilyId() uuid.UUID
	MemberId() uuid.UUID

	Equal(other Payer) bool
}

func NewPayer(
	payerType PayerType,
	familyId uuid.UUID,
	memberId *uuid.UUID) Payer {
	switch payerType {
	case FamilyPayer:
		return NewFamilyPayer(familyId)
	case FamilyMemberPayer:
		if memberId == nil {
			panic("member id should not be nil")
		}
		return NewFamilyMemberPayer(familyId, *memberId)
	default:
		panic("unknown payer type")
	}
}

type familyPayer struct {
	familyId uuid.UUID
}

func NewFamilyPayer(
	familyId uuid.UUID) Payer {
	return familyPayer{
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

func (f familyPayer) ETagFields() []interface{} {
	return []interface{}{
		f.familyId,
		FamilyMemberPayer.String(),
	}

}
func (f familyPayer) ETag() string {
	return entity.CalculateETag(f)
}

func (f familyPayer) Equal(other Payer) bool {
	if other == nil {
		return false
	}

	return f.ETag() == other.ETag()
}

type familyMemberPayer struct {
	*entity.Base

	familyId uuid.UUID
	memberId uuid.UUID
}

func NewFamilyMemberPayer(
	familyId uuid.UUID,
	memberId uuid.UUID) Payer {
	return familyMemberPayer{
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

func (f familyMemberPayer) ETagFields() []interface{} {
	return []interface{}{
		f.familyId,
		f.memberId,
		FamilyMemberPayer.String(),
	}

}
func (f familyMemberPayer) ETag() string {
	return entity.CalculateETag(f)
}

func (f familyMemberPayer) Equal(other Payer) bool {
	if other == nil {
		return false
	}

	return f.ETag() == other.ETag()
}

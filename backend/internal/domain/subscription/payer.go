package subscription

import (
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type PayerType string

const (
	UnknownPayer      PayerType = "unknown"
	PersonalPayer     PayerType = "personal"
	FamilyMemberPayer PayerType = "family_member"
	FamilyPayer       PayerType = "family"
)

func (p PayerType) String() string {
	return string(p)
}

func ParsePayerType(input string) (payerType PayerType, err error) {
	switch input {
	case string(PersonalPayer):
		return PersonalPayer, nil
	case string(FamilyMemberPayer):
		return FamilyMemberPayer, nil
	case string(FamilyPayer):
		return FamilyPayer, nil
	default:
		return UnknownPayer, ErrUnknownPayerType
	}
}

func MustParsePayerType(input string) PayerType {
	t, err := ParsePayerType(input)
	if err != nil {
		panic(err)
	}
	return t
}

type Payer interface {
	entity.ETagEntity

	Type() PayerType
	FamilyId() types.FamilyID
	MemberId() types.FamilyMemberID

	Equal(other Payer) bool
}

func NewPayer(
	payerType PayerType,
	familyId types.FamilyID,
	memberId *types.FamilyMemberID) Payer {
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

func NewPayerFromOwner(owner types.Owner,
	payerType PayerType,
	familyMemberId *types.FamilyMemberID) (Payer, error) {
	if owner.Type() != types.FamilyOwnerType {
		return nil, ErrPayerNeedFamily
	}

	return NewPayer(payerType, owner.FamilyId(), familyMemberId), nil
}

type familyPayer struct {
	familyId types.FamilyID
}

func NewFamilyPayer(
	familyId types.FamilyID) Payer {
	return familyPayer{
		familyId: familyId,
	}
}

func (f familyPayer) Type() PayerType {
	return FamilyPayer
}

func (f familyPayer) FamilyId() types.FamilyID {
	return f.familyId
}

func (f familyPayer) MemberId() types.FamilyMemberID {
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
	familyId types.FamilyID
	memberId types.FamilyMemberID
}

func NewFamilyMemberPayer(
	familyId types.FamilyID,
	memberId types.FamilyMemberID) Payer {
	return familyMemberPayer{
		familyId: familyId,
		memberId: memberId,
	}
}

func (f familyMemberPayer) FamilyId() types.FamilyID {
	return f.familyId
}

func (f familyMemberPayer) MemberId() types.FamilyMemberID {
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

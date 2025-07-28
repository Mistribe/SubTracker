package subscription

import (
	"github.com/google/uuid"
)

type PayerType string

const (
	FamilyMemberPayer PayerType = "family_member"
	FamilyPayer       PayerType = "family"
)

type Payer interface {
	Id() uuid.UUID
	Type() PayerType
}

type familyPayer struct {
	familyId uuid.UUID
}

func NewFamilyPayer(familyId uuid.UUID) Payer {
	return familyPayer{familyId: familyId}
}

func (f familyPayer) Id() uuid.UUID {
	return f.familyId
}

func (f familyPayer) Type() PayerType {
	return FamilyPayer
}

type familyMemberPayer struct {
	familyId uuid.UUID
	memberId uuid.UUID
}

func NewFamilyMemberPayer(familyId uuid.UUID, memberId uuid.UUID) Payer {
	return familyMemberPayer{familyId: familyId, memberId: memberId}
}

func (f familyMemberPayer) Id() uuid.UUID {
	return f.memberId
}

func (f familyMemberPayer) Type() PayerType {
	return FamilyMemberPayer
}

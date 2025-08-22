package family

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/validationx"
)

type MemberType string

const (
	UnknownMemberType MemberType = "unknown"
	OwnerMemberType   MemberType = "owner"
	AdultMemberType   MemberType = "adult"
	KidMemberType     MemberType = "kid"
)

func (t MemberType) String() string {
	return string(t)
}

func MustParseMemberType(input string) MemberType {
	switch input {
	case string(OwnerMemberType):
		return OwnerMemberType
	case string(AdultMemberType):
		return AdultMemberType
	case string(KidMemberType):
		return KidMemberType
	default:
		panic("unknown member type")
	}
}

func ParseMemberType(input string) (MemberType, error) {
	switch input {
	case string(OwnerMemberType):
		return OwnerMemberType, nil
	case string(AdultMemberType):
		return AdultMemberType, nil
	case string(KidMemberType):
		return KidMemberType, nil
	default:
		return UnknownMemberType, errors.New("unknown member type")
	}

}

type Member interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	Type() MemberType
	UserId() *string
	SetName(string)
	SetType(MemberType)
	SetUserId(*string)
	FamilyId() uuid.UUID

	Equal(member Member) bool
	GetValidationErrors() validationx.Errors
	SetInvitationCode(code string) bool
	InvitationCode() *string
}

type member struct {
	*entity.Base

	familyId       uuid.UUID
	name           string
	userId         *string
	userType       MemberType
	invitationCode *string
}

func NewMember(
	id uuid.UUID,
	familyId uuid.UUID,
	name string,
	userType MemberType,
	invitationCode *string,
	createdAt time.Time,
	updatedAt time.Time) Member {
	return &member{
		Base:           entity.NewBase(id, createdAt, updatedAt, true, false),
		familyId:       familyId,
		name:           strings.TrimSpace(name),
		userId:         nil,
		userType:       userType,
		invitationCode: invitationCode,
	}
}

func (m *member) SetInvitationCode(code string) bool {
	m.invitationCode = &code
	m.SetAsDirty()
	return true
}

func (m *member) InvitationCode() *string {
	return m.invitationCode
}

func (m *member) Name() string {
	return m.name
}

func (m *member) Type() MemberType {
	return m.userType
}

func (m *member) UserId() *string {
	return m.userId
}

func (m *member) SetUserId(userId *string) {
	if m.userId == userId {
		return
	}
	m.userId = userId
	m.SetAsDirty()
}

func (m *member) ETagFields() []interface{} {
	return []interface{}{
		m.familyId.String(),
		m.name,
		m.userId,
		m.userType.String(),
	}
}
func (m *member) ETag() string {
	return entity.CalculateETag(m)
}

func (m *member) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if m.name == "" {
		errors = append(errors, validationx.NewError("name", "name is empty"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

func (m *member) SetName(name string) {
	m.name = name
	m.SetAsDirty()
}

func (m *member) SetType(userType MemberType) {
	m.userType = userType
	m.SetAsDirty()
}

func (m *member) FamilyId() uuid.UUID {
	return m.familyId
}

func (m *member) Equal(other Member) bool {
	if other == nil {
		return false
	}

	return m.ETag() == other.ETag()
}

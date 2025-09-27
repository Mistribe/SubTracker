package family

import (
	"errors"
	"strings"
	"time"

	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x/validation"
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
	entity.Entity[types.FamilyMemberID]
	entity.ETagEntity

	Name() string
	Type() MemberType
	UserId() *types.UserID
	FamilyId() types.FamilyID
	SetUserId(*types.UserID)
	SetName(string)
	SetType(MemberType)

	Equal(member Member) bool
	GetValidationErrors() validation.Errors
	SetInvitationCode(code *string) bool
	InvitationCode() *string
}

type member struct {
	*entity.Base[types.FamilyMemberID]

	familyId       types.FamilyID
	name           string
	userId         *types.UserID
	userType       MemberType
	invitationCode *string
}

func NewMember(
	id types.FamilyMemberID,
	familyId types.FamilyID,
	name string,
	userType MemberType,
	invitationCode *string,
	createdAt time.Time,
	updatedAt time.Time) Member {
	return &member{
		Base:           entity.NewBase[types.FamilyMemberID](id, createdAt, updatedAt, true, false),
		familyId:       familyId,
		name:           strings.TrimSpace(name),
		userId:         nil,
		userType:       userType,
		invitationCode: invitationCode,
	}
}

func (m *member) SetInvitationCode(code *string) bool {
	m.invitationCode = code
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

func (m *member) UserId() *types.UserID {
	return m.userId
}

func (m *member) SetUserId(userId *types.UserID) {
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

func (m *member) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if m.name == "" {
		errors = append(errors, validation.NewError("name", "name is empty"))
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

func (m *member) FamilyId() types.FamilyID {
	return m.familyId
}

func (m *member) Equal(other Member) bool {
	if other == nil {
		return false
	}

	return m.ETag() == other.ETag()
}

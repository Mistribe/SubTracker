package family

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Member struct {
	id        uuid.UUID
	familyId  uuid.UUID
	name      string
	email     *string
	userId    *string
	isKid     bool
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
	isExists  bool
}

func NewMember(
	id uuid.UUID,
	familyId uuid.UUID,
	name string,
	email option.Option[string],
	isKid bool,
	createdAt time.Time,
	updatedAt time.Time) result.Result[Member] {
	mbr := NewMemberWithoutValidation(id,
		familyId,
		strings.TrimSpace(name),
		email.Value(),
		isKid,
		createdAt,
		updatedAt,
		false)

	if err := mbr.Validate(); err != nil {
		return result.Fail[Member](err)
	}

	return result.Success(mbr)
}

func NewMemberWithoutValidation(
	id uuid.UUID,
	familyId uuid.UUID,
	name string,
	email *string,
	isKid bool,
	createdAt time.Time,
	updatedAt time.Time,
	isExists bool) Member {
	if email != nil {
		trimEmail := strings.TrimSpace(*email)
		email = &trimEmail
	}
	return Member{
		id:        id,
		familyId:  familyId,
		name:      strings.TrimSpace(name),
		email:     email,
		userId:    nil,
		isKid:     isKid,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   true,
		isExists:  isExists,
	}
}

func (m *Member) Id() uuid.UUID {
	return m.id
}

func (m *Member) Name() string {
	return m.name
}

func (m *Member) IsKid() bool {
	return m.isKid
}

func (m *Member) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Member) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Member) Validate() error {
	if m.email != nil {
		if !emailRegex.MatchString(*m.email) {
			return errors.New("invalid email format")
		}
	}
	if m.name == "" {
		return errors.New("name is empty")
	}
	return nil
}

func (m *Member) SetName(name string) {
	m.name = name
	m.isDirty = true
}

func (m *Member) SetAsKid() {
	m.isKid = true
	m.isDirty = true
}

func (m *Member) SetAsAdult() {
	m.isKid = false
	m.isDirty = true
}

func (m *Member) SetUpdatedAt(updatedAt time.Time) {
	m.updatedAt = updatedAt
	m.isDirty = true
}

func (m *Member) Email() option.Option[string] {
	return option.New(m.email)
}

func (m *Member) SetEmail(email option.Option[string]) {
	m.email = email.Value()
	m.isDirty = true
}

func (m *Member) Clean() {
	m.isDirty = false
	m.isExists = true
}

func (m *Member) IsDirty() bool {
	return m.isDirty
}

func (m *Member) IsExists() bool {
	return m.isExists
}

func (m *Member) FamilyId() uuid.UUID {
	return m.familyId
}

func (m *Member) Equal(member Member) bool {
	if m.id != member.id ||
		m.familyId != member.familyId ||
		m.name != member.name ||
		m.isKid != member.isKid ||
		!m.createdAt.Equal(member.createdAt) ||
		!m.updatedAt.Equal(member.updatedAt) {
		return false
	}

	if (m.email == nil && member.email != nil) ||
		(m.email != nil && member.email == nil) {
		return false
	}

	if m.email != nil && member.email != nil && *m.email != *member.email {
		return false
	}

	if (m.userId == nil && member.userId != nil) ||
		(m.userId != nil && member.userId == nil) {
		return false
	}

	if m.userId != nil && member.userId != nil && *m.userId != *member.userId {
		return false
	}

	return true
}

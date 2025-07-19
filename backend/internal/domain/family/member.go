package family

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Member struct {
	*entity.Base

	familyId uuid.UUID
	name     string
	email    *string
	userId   *string
	isKid    bool
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
		Base:     entity.NewBase(id, createdAt, updatedAt, true, isExists),
		familyId: familyId,
		name:     strings.TrimSpace(name),
		email:    email,
		userId:   nil,
		isKid:    isKid,
	}
}

func (m *Member) Name() string {
	return m.name
}

func (m *Member) IsKid() bool {
	return m.isKid
}

func (m *Member) ETagFields() []interface{} {
	return []interface{}{
		m.familyId.String(),
		m.name,
		m.email,
		m.userId,
		m.isKid,
	}
}
func (m *Member) ETag() string {
	return entity.CalculateETag(m, m.Base)
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
	m.SetAsDirty()
}

func (m *Member) SetAsKid() {
	m.isKid = true
	m.SetAsDirty()
}

func (m *Member) SetAsAdult() {
	m.isKid = false
	m.SetAsDirty()
}

func (m *Member) Email() option.Option[string] {
	return option.New(m.email)
}

func (m *Member) SetEmail(email option.Option[string]) {
	m.email = email.Value()
	m.SetAsDirty()
}

func (m *Member) FamilyId() uuid.UUID {
	return m.familyId
}

func (m *Member) Equal(member Member) bool {
	if !m.Base.Equal(*member.Base) ||
		m.familyId != member.familyId ||
		m.name != member.name ||
		m.isKid != member.isKid {
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

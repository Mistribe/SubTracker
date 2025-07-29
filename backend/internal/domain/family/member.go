package family

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/validationx"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Member interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	IsKid() bool
	SetName(string)
	SetAsKid()
	SetAsAdult()
	SetUserId(*string)
	FamilyId() uuid.UUID
	Equal(member Member) bool
	GetValidationErrors() validationx.Errors
}

type member struct {
	*entity.Base

	familyId uuid.UUID
	name     string
	userId   *string
	isKid    bool
}

func NewMember(
	id uuid.UUID,
	familyId uuid.UUID,
	name string,
	isKid bool,
	createdAt time.Time,
	updatedAt time.Time) Member {
	return &member{
		Base:     entity.NewBase(id, createdAt, updatedAt, true, false),
		familyId: familyId,
		name:     strings.TrimSpace(name),
		userId:   nil,
		isKid:    isKid,
	}
}

func (m *member) Name() string {
	return m.name
}

func (m *member) IsKid() bool {
	return m.isKid
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
		m.isKid,
	}
}
func (m *member) ETag() string {
	return entity.CalculateETag(m, m.Base)
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

func (m *member) SetAsKid() {
	m.isKid = true
	m.SetAsDirty()
}

func (m *member) SetAsAdult() {
	m.isKid = false
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

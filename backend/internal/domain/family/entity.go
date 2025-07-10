package family

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type Member struct {
	id        uuid.UUID
	name      string
	isKid     bool
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
}

func NewMember(id uuid.UUID, name string, isKid bool, createdAt time.Time, updatedAt time.Time) result.Result[Member] {
	mbr := NewMemberWithoutValidation(id, name, isKid, createdAt, updatedAt)

	if err := mbr.Validate(); err != nil {
		return result.Fail[Member](err)
	}

	return result.Success(mbr)
}

func NewMemberWithoutValidation(
	id uuid.UUID,
	name string,
	isKid bool,
	createdAt time.Time,
	updatedAt time.Time) Member {
	return Member{
		id:        id,
		name:      name,
		isKid:     isKid,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   true,
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

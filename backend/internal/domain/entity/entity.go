package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entity interface {
	Id() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Clean()
	IsDirty() bool
	SetUpdatedAt(updatedAt time.Time)
	SetAsDirty()
	IsExists() bool
	Equal(other Base) bool
}

type ETagEntity interface {
	ETagFields() []interface{}
	ETag() string
}

type Base struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
	isExists  bool
}

func NewBase(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	isDirty bool,
	isExists bool) *Base {
	return &Base{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   isDirty,
		isExists:  isExists,
	}
}

func (b *Base) Id() uuid.UUID {
	return b.id
}

func (b *Base) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Base) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Base) Clean() {
	b.isDirty = false
	b.isExists = true
}

func (b *Base) IsDirty() bool {
	return b.isDirty
}

func (b *Base) SetUpdatedAt(updatedAt time.Time) {
	b.updatedAt = updatedAt
	b.SetAsDirty()
}

func (b *Base) SetAsDirty() {
	b.isDirty = true
}

func (b *Base) IsExists() bool {
	return b.isExists
}

func (b *Base) Equal(other Base) bool {
	return b.id == other.id &&
		b.createdAt == other.createdAt &&
		b.updatedAt == other.updatedAt
}

func (b *Base) ETagFields() []interface{} {
	return []interface{}{
		b.id.String(),
		b.createdAt,
		b.updatedAt,
	}
}

func (b *Base) ETag() string {
	return CalculateETag(b)
}

package entity

import (
	"time"
)

type Entity[TKey comparable] interface {
	Id() TKey
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Clean()
	IsDirty() bool
	SetUpdatedAt(updatedAt time.Time)
	SetAsDirty()
	IsExists() bool
}

type Base[TKey comparable] struct {
	id        TKey
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
	isExists  bool
}

func NewBase[TKey comparable](
	id TKey,
	createdAt time.Time,
	updatedAt time.Time,
	isDirty bool,
	isExists bool) *Base[TKey] {
	return &Base[TKey]{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   isDirty,
		isExists:  isExists,
	}
}

func (b *Base[TKey]) Id() TKey {
	return b.id
}

func (b *Base[TKey]) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Base[TKey]) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Base[TKey]) Clean() {
	b.isDirty = false
	b.isExists = true
}

func (b *Base[TKey]) IsDirty() bool {
	return b.isDirty
}

func (b *Base[TKey]) SetUpdatedAt(updatedAt time.Time) {
	b.updatedAt = updatedAt
	b.SetAsDirty()
}

func (b *Base[TKey]) SetAsDirty() {
	b.isDirty = true
}

func (b *Base[TKey]) IsExists() bool {
	return b.isExists
}

func (b *Base[TKey]) Equal(other Base[TKey]) bool {
	return b.id == other.id &&
		b.createdAt == other.createdAt &&
		b.updatedAt == other.updatedAt
}

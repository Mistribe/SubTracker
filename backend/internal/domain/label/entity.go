package label

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type Label struct {
	id        uuid.UUID
	name      string
	isDefault bool
	color     string
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
}

func NewLabel(id uuid.UUID, name string, isDefault bool, color string, createdAt time.Time,
	updatedAt time.Time) result.Result[Label] {
	lbl := Label{
		id:        id,
		name:      name,
		isDefault: isDefault,
		color:     color,
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   true,
	}

	if err := lbl.Validate(); err != nil {
		return result.Fail[Label](err)
	}

	return result.Success(lbl)
}

func (l *Label) Id() uuid.UUID {
	return l.id
}

func (l *Label) Name() string {
	return l.name
}

func (l *Label) IsDefault() bool {
	return l.isDefault
}

func (l *Label) Color() string {
	return l.color
}

func (l *Label) CreatedAt() time.Time {
	return l.createdAt
}

func (l *Label) UpdatedAt() time.Time {
	return l.updatedAt
}

func (l *Label) Validate() error {
	// todo
	return nil
}

func (l *Label) SetName(name string) {
	l.name = name
	l.isDirty = true
}

func (l *Label) SetIsDefault(isDefault bool) {
	l.isDefault = isDefault
	l.isDirty = true
}

func (l *Label) SetColor(color string) {
	l.color = color
	l.isDirty = true
}

func (l *Label) SetUpdatedAt(updatedAt time.Time) {
	l.updatedAt = updatedAt
	l.isDirty = true
}

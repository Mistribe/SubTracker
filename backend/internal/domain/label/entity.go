package label

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/result"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{8}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

type Label struct {
	id        uuid.UUID
	name      string
	isDefault bool
	color     string
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
	isExists  bool
}

func NewLabelWithoutValidation(
	id uuid.UUID,
	name string,
	isDefault bool,
	color string,
	createdAt time.Time,
	updatedAt time.Time,
	isExists bool) Label {
	return Label{
		id:        id,
		name:      strings.TrimSpace(name),
		isDefault: isDefault,
		color:     strings.TrimSpace(color),
		createdAt: createdAt,
		updatedAt: updatedAt,
		isDirty:   true,
		isExists:  isExists,
	}
}

func NewLabel(
	id uuid.UUID,
	name string,
	isDefault bool,
	color string,
	createdAt time.Time,
	updatedAt time.Time) result.Result[Label] {
	lbl := NewLabelWithoutValidation(id, name, isDefault, color, createdAt, updatedAt, false)

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
	if err := l.validateName(); err != nil {
		return err
	}

	if err := l.validateColor(); err != nil {
		return err
	}

	return nil
}

func (l *Label) validateName() error {
	name := strings.TrimSpace(l.name)
	if name == "" {
		return ErrLabelNameEmpty
	}

	if len(name) > 100 {
		return ErrLabelNameTooLong
	}

	return nil
}

func (l *Label) validateColor() error {
	color := strings.TrimSpace(l.color)
	if color == "" || !hexColorRegex.MatchString(color) {
		return ErrLabelColorInvalid
	}

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

func (l *Label) IsExists() bool {
	return l.isExists
}

func (l *Label) Equal(other Label) bool {
	return l.id == other.id &&
		l.name == other.name &&
		l.isDefault == other.isDefault &&
		l.color == other.color &&
		l.createdAt.Equal(other.createdAt) &&
		l.updatedAt.Equal(other.updatedAt)
}

func (l *Label) Clean() {
	l.isDirty = false
	l.isExists = true
}

func (l *Label) IsDirty() bool {
	return l.isDirty
}

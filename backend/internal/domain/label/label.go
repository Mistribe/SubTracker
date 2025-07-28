package label

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/validationx"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{8}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

type Label interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	Owner() user.Owner
	IsDefault() bool
	Color() string
	GetValidationErrors() validationx.Errors
	SetName(name string)
	SetIsDefault(isDefault bool)
	SetColor(color string)
	Equal(other Label) bool
}

type label struct {
	*entity.Base

	owner     user.Owner
	name      string
	isDefault bool
	color     string
}

func NewLabel(
	id uuid.UUID,
	owner user.Owner,
	name string,
	isDefault bool,
	color string,
	createdAt time.Time,
	updatedAt time.Time) Label {
	return &label{
		Base:      entity.NewBase(id, createdAt, updatedAt, true, false),
		owner:     owner,
		name:      strings.TrimSpace(name),
		isDefault: isDefault,
		color:     strings.TrimSpace(color),
	}
}

func (l *label) Name() string {
	return l.name
}

func (l *label) Owner() user.Owner {
	return l.owner
}

func (l *label) IsDefault() bool {
	return l.isDefault
}

func (l *label) Color() string {
	return l.color
}

func (l *label) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	name := strings.TrimSpace(l.name)
	if name == "" {
		errors = append(errors, validationx.NewError("name", "name is required and cannot be empty"))
	}

	if len(name) > 100 {
		errors = append(errors, validationx.NewError("name", "name cannot be longer than 100 characters"))
	}

	color := strings.TrimSpace(l.color)
	if color == "" || !hexColorRegex.MatchString(color) {
		errors = append(errors, validationx.NewError("color", "color is required and must be a valid hex color"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

func (l *label) SetName(name string) {
	l.name = name
	l.SetAsDirty()
}

func (l *label) SetIsDefault(isDefault bool) {
	l.isDefault = isDefault
	l.SetAsDirty()
}

func (l *label) SetColor(color string) {
	l.color = color
	l.SetAsDirty()
}

func (l *label) Equal(other Label) bool {
	if other == nil {
		return false
	}

	return l.ETag() == other.ETag()
}

func (l *label) ETagFields() []interface{} {
	return []interface{}{
		l.name,
		l.isDefault,
		l.color,
	}
}
func (l *label) ETag() string {
	return entity.CalculateETag(l, l.Base)
}

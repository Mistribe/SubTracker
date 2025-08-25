package label

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/x/validation"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{8}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

type Label interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	Key() *string
	Owner() auth.Owner
	Color() string
	GetValidationErrors() validation.Errors
	SetName(name string)
	SetColor(color string)
	Equal(other Label) bool
}

type label struct {
	*entity.Base

	owner     auth.Owner
	name      string
	key       *string
	isDefault bool
	color     string
}

func NewLabel(
	id uuid.UUID,
	owner auth.Owner,
	name string,
	key *string,
	color string,
	createdAt time.Time,
	updatedAt time.Time) Label {
	return &label{
		Base:  entity.NewBase(id, createdAt, updatedAt, true, false),
		owner: owner,
		name:  strings.TrimSpace(name),
		key:   key,
		color: strings.TrimSpace(color),
	}
}

func (l *label) Name() string {
	return l.name
}

func (l *label) Key() *string {
	return l.key
}

func (l *label) Owner() auth.Owner {
	return l.owner
}

func (l *label) Color() string {
	return l.color
}

func (l *label) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	name := strings.TrimSpace(l.name)
	if name == "" {
		errors = append(errors, validation.NewError("name", "name is required and cannot be empty"))
	}

	if len(name) > 100 {
		errors = append(errors, validation.NewError("name", "name cannot be longer than 100 characters"))
	}

	color := strings.TrimSpace(l.color)
	if color == "" || !hexColorRegex.MatchString(color) {
		errors = append(errors, validation.NewError("color", "color is required and must be a valid hex color"))
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

func (l *label) SetColor(color string) {
	if l.color == color {
		return
	}
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
	return entity.CalculateETag(l)
}

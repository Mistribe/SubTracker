package label

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{8}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

type Label struct {
	*entity.Base

	ownerId   *string
	name      string
	isDefault bool
	color     string
}

func NewLabelWithoutValidation(
	id uuid.UUID,
	ownerId *string,
	name string,
	isDefault bool,
	color string,
	createdAt time.Time,
	updatedAt time.Time,
	isExists bool) Label {
	return Label{
		Base:      entity.NewBase(id, createdAt, updatedAt, true, isExists),
		ownerId:   ownerId,
		name:      strings.TrimSpace(name),
		isDefault: isDefault,
		color:     strings.TrimSpace(color),
	}
}

func NewLabel(
	id uuid.UUID,
	ownerId *string,
	name string,
	isDefault bool,
	color string,
	createdAt time.Time,
	updatedAt time.Time) result.Result[Label] {
	lbl := NewLabelWithoutValidation(id, ownerId, name, isDefault, color, createdAt, updatedAt, false)

	if err := lbl.Validate(); err != nil {
		return result.Fail[Label](err)
	}

	return result.Success(lbl)
}

func (l *Label) Name() string {
	return l.name
}

func (l *Label) OwnerId() *string {
	return l.ownerId
}

func (l *Label) IsDefault() bool {
	return l.isDefault
}

func (l *Label) Color() string {
	return l.color
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
	l.SetAsDirty()
}

func (l *Label) SetIsDefault(isDefault bool) {
	l.isDefault = isDefault
	l.SetAsDirty()
}

func (l *Label) SetColor(color string) {
	l.color = color
	l.SetAsDirty()
}

func (l *Label) Equal(other Label) bool {
	return l.Base.Equal(*other.Base) &&
		l.name == other.name &&
		l.isDefault == other.isDefault &&
		l.color == other.color
}

func (l *Label) ETagFields() []interface{} {
	return []interface{}{
		l.name,
		l.isDefault,
		l.color,
	}
}
func (l *Label) ETag() string {
	return entity.CalculateETag(l, l.Base)
}

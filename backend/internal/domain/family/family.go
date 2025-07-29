package family

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"github.com/oleexo/subtracker/pkg/validationx"
)

type Family interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	SetName(name string)
	OwnerId() string
	Members() *slicesx.Tracked[Member]
	AddMember(member Member) error
	GetMember(id uuid.UUID) Member
	UpdateMember(member Member) error
	ContainsMember(id uuid.UUID) bool
	Equal(family Family) bool
	GetValidationErrors() validationx.Errors
}
type family struct {
	*entity.Base

	ownerId string
	name    string
	members *slicesx.Tracked[Member]
}

func NewFamily(
	id uuid.UUID,
	ownerId string,
	name string,
	members []Member,
	createdAt time.Time,
	updatedAt time.Time) Family {
	return &family{
		Base:    entity.NewBase(id, createdAt, updatedAt, true, false),
		ownerId: ownerId,
		name:    name,
		members: slicesx.NewTracked(members, memberUniqueComparer, memberComparer),
	}
}

func (f *family) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if f.name == "" {
		errors = append(errors, validationx.NewError("name", "name is empty"))
	}
	if len(f.name) < 3 {
		errors = append(errors, validationx.NewError("name", "name must be at least 3 characters long"))
	}
	if len(f.name) > 100 {
		errors = append(errors, validationx.NewError("name", "name cannot be longer than 100 characters"))
	}

	idMap := make(map[uuid.UUID]bool)

	for mbr := range f.members.It() {
		if idMap[mbr.Id()] {
			errors = append(errors, validationx.NewError("members", "duplicate member id"))
		}
		idMap[mbr.Id()] = true

		errors = append(errors, mbr.GetValidationErrors()...)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

func (f *family) Name() string {
	return f.name
}

func (f *family) OwnerId() string {
	return f.ownerId
}

func (f *family) Members() *slicesx.Tracked[Member] {
	return f.members
}

func (f *family) AddMember(member Member) error {
	if f.isDuplicateMember(member) {
		return ErrDuplicateMember
	}

	f.members.Add(member)
	f.SetAsDirty()
	return nil
}

func (f *family) isDuplicateMember(member Member) bool {
	for m := range f.members.It() {
		if m.Id() == member.Id() {
			return true
		}
	}
	return false
}

func (f *family) GetMember(id uuid.UUID) Member {
	for m := range f.members.It() {
		if m.Id() == id {
			return m
		}
	}

	return nil
}

func (f *family) UpdateMember(member Member) error {
	if !f.members.Update(member) {
		return ErrFamilyMemberNotFound

	}

	return nil
}

func (f *family) ContainsMember(id uuid.UUID) bool {
	for m := range f.members.It() {
		if m.Id() == id {
			return true
		}
	}

	return false
}

func (f *family) SetName(name string) {
	f.name = name
	f.SetAsDirty()
}

func (f *family) Equal(other Family) bool {
	if other == nil {
		return false
	}

	return f.ETag() == other.ETag()
}

func (f *family) ETagFields() []interface{} {
	fields := []interface{}{
		f.ownerId,
		f.name,
	}

	for mbr := range f.members.It() {
		fields = append(fields, mbr.ETagFields()...)
	}

	return fields
}
func (f *family) ETag() string {
	return entity.CalculateETag(f, f.Base)
}

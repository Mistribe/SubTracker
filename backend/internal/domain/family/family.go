package family

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/herd"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type Family interface {
	entity.Entity[types.FamilyID]
	entity.ETagEntity

	Name() string
	SetName(name string)
	Owner() types.Owner
	Members() *slicesx.Tracked[Member]
	AddMember(member Member) error
	RemoveMember(member Member) error
	GetMember(id types.FamilyMemberID) Member
	UpdateMember(member Member) error
	ContainsMember(id types.FamilyMemberID) bool
	Equal(family Family) bool
	GetValidationErrors() validation.Errors
}
type family struct {
	*entity.Base[types.FamilyID]

	owner   types.Owner
	name    string
	members *slicesx.Tracked[Member]
}

func NewFamily(
	id types.FamilyID,
	ownerId types.UserID,
	name string,
	members []Member,
	createdAt time.Time,
	updatedAt time.Time) Family {
	return &family{
		Base:    entity.NewBase[types.FamilyID](id, createdAt, updatedAt, true, false),
		owner:   types.NewPersonalOwner(ownerId),
		name:    name,
		members: slicesx.NewTracked(members, memberUniqueComparer, memberComparer),
	}
}

func (f *family) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if f.name == "" {
		errors = append(errors, validation.NewError("name", "name is empty"))
	}
	if len(f.name) < 3 {
		errors = append(errors, validation.NewError("name", "name must be at least 3 characters long"))
	}
	if len(f.name) > 100 {
		errors = append(errors, validation.NewError("name", "name cannot be longer than 100 characters"))
	}

	idMap := herd.NewSet[types.FamilyMemberID]()

	if f.members.Len() == 0 {
		errors = append(errors, validation.NewError("members", "family must have at least one member"))
	}

	for mbr := range f.members.It() {
		if idMap.Contains(mbr.Id()) {
			errors = append(errors, validation.NewError("members", "duplicate member id"))
		}
		idMap.Add(mbr.Id())

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

func (f *family) Owner() types.Owner {
	return f.owner
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

func (f *family) RemoveMember(member Member) error {
	if !f.members.Remove(member) {
		return ErrFamilyMemberNotFound
	}

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

func (f *family) GetMember(id types.FamilyMemberID) Member {
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

	f.SetAsDirty()
	return nil
}

func (f *family) ContainsMember(id types.FamilyMemberID) bool {
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
		f.name,
	}

	for mbr := range f.members.It() {
		fields = append(fields, mbr.ETagFields()...)
	}

	fields = append(fields, f.owner.ETag())

	return fields
}
func (f *family) ETag() string {
	return entity.CalculateETag(f)
}

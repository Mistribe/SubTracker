package family

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Family struct {
	*entity.Base

	ownerId          string
	name             string
	members          []Member
	haveJointAccount bool
}

var ErrDuplicateMember = errors.New("member with the same ID or email already exists")

func NewFamily(
	id uuid.UUID,
	ownerId string,
	name string,
	haveJointAccount bool,
	createdAt time.Time,
	updatedAt time.Time) result.Result[Family] {
	f := NewFamilyWithoutValidation(id, ownerId, name, haveJointAccount, nil, createdAt, updatedAt, false)

	if err := f.Validate(); err != nil {
		return result.Fail[Family](err)
	}
	return result.Success(f)
}

func NewFamilyWithoutValidation(
	id uuid.UUID,
	ownerId string,
	name string,
	haveJointAccount bool,
	members []Member,
	createdAt time.Time,
	updatedAt time.Time,
	isExists bool) Family {
	return Family{
		Base:             entity.NewBase(id, createdAt, updatedAt, true, isExists),
		ownerId:          ownerId,
		name:             name,
		members:          members,
		haveJointAccount: haveJointAccount,
	}
}

func (f *Family) Validate() error {
	if f.name == "" {
		return ErrNameRequired
	}
	if len(f.name) < 3 {
		return ErrNameTooShort
	}
	if len(f.name) > 100 {
		return ErrNameTooLong
	}

	emailMap := make(map[string]bool)
	idMap := make(map[uuid.UUID]bool)

	for _, member := range f.members {
		if idMap[member.Id()] {
			return ErrDuplicateMember
		}
		idMap[member.Id()] = true

		if err := option.Match(member.Email(), func(email string) error {
			if emailMap[email] {
				return ErrDuplicateMember
			}
			emailMap[email] = true
			return nil
		}, func() error {
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Family) HaveJointAccount() bool {
	return f.haveJointAccount
}

func (f *Family) Name() string {
	return f.name
}

func (f *Family) OwnerId() string {
	return f.ownerId
}

func (f *Family) Members() []Member {
	return f.members
}

func (f *Family) AddMember(member Member) error {
	if f.isDuplicateMember(member) {
		return ErrDuplicateMember
	}

	f.members = append(f.members, member)
	f.SetAsDirty()
	return nil
}

func (f *Family) isDuplicateMember(member Member) bool {
	for _, m := range f.members {
		if m.Id() == member.Id() {
			return true
		}
		if m.Email().IsSome() && member.Email().IsSome() && m.Email().Value() == member.Email().Value() {
			return true
		}
	}
	return false
}

func (f *Family) GetMember(id uuid.UUID) option.Option[Member] {
	for _, m := range f.members {
		if m.Id() == id {
			return option.Some(m)
		}
	}

	return option.None[Member]()
}

func (f *Family) UpdateMember(member Member) error {
	for idx, m := range f.members {
		if m.Id() == member.Base.Id() {
			f.members[idx] = member
			return nil
		}
	}

	return ErrFamilyMemberNotFound
}

func (f *Family) ContainsMember(id uuid.UUID) bool {
	for _, m := range f.members {
		if m.Id() == id {
			return true
		}
	}

	return false
}

func (f *Family) SetHaveJointAccount(haveJointAccount bool) {
	f.haveJointAccount = haveJointAccount
	f.SetAsDirty()
}

func (f *Family) SetName(name string) {
	f.name = name
	f.SetAsDirty()
}

func (f *Family) Equal(family Family) bool {
	if !f.Base.Equal(*family.Base) ||
		f.ownerId != family.ownerId ||
		f.name != family.name ||
		f.haveJointAccount != family.haveJointAccount ||
		len(f.members) != len(family.members) {
		return false
	}

	memberMap := make(map[uuid.UUID]Member)
	for _, member := range f.members {
		memberMap[member.Id()] = member
	}

	for _, member := range family.members {
		if existingMember, ok := memberMap[member.Id()]; !ok || !existingMember.Equal(member) {
			return false
		}
	}

	return true
}

func (f *Family) ETagFields() []interface{} {
	fields := []interface{}{
		f.ownerId,
		f.haveJointAccount,
		f.name,
	}

	for _, member := range f.members {
		fields = append(fields, member.ETagFields()...)
	}

	return fields
}
func (f *Family) ETag() string {
	return entity.CalculateETag(f, f.Base)
}

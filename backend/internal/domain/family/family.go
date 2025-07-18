package family

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Family struct {
	id               uuid.UUID
	ownerId          string
	name             string
	members          []Member
	haveJointAccount bool
	createdAt        time.Time
	updatedAt        time.Time
	isDirty          bool
	isExists         bool
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
		id:               id,
		ownerId:          ownerId,
		name:             name,
		members:          members,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		haveJointAccount: haveJointAccount,
		isDirty:          true,
		isExists:         isExists,
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

func (f *Family) Id() uuid.UUID {
	return f.id
}
func (f *Family) HaveJointAccount() bool {
	return f.haveJointAccount
}
func (f *Family) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Family) UpdatedAt() time.Time {
	return f.updatedAt
}

func (f *Family) Name() string {
	return f.name
}

func (f *Family) OwnerId() string {
	return f.ownerId
}

func (f *Family) Clean() {
	f.isDirty = false
	f.isExists = true
}

func (f *Family) IsDirty() bool {
	return f.isDirty
}

func (f *Family) IsExists() bool {
	return f.isExists
}

func (f *Family) Members() []Member {
	return f.members
}

func (f *Family) AddMember(member Member) error {
	if f.isDuplicateMember(member) {
		return ErrDuplicateMember
	}

	f.members = append(f.members, member)
	f.isDirty = true
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
		if m.Id() == member.id {
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

func (f *Family) SetUpdatedAt(updatedAt time.Time) {
	f.updatedAt = updatedAt
	f.isDirty = true
}

func (f *Family) SetHaveJointAccount(haveJointAccount bool) {
	f.haveJointAccount = haveJointAccount
	f.isDirty = true
}

func (f *Family) SetName(name string) {
	f.name = name
	f.isDirty = true
}

func (f *Family) Equal(family Family) bool {
	if f.id != family.id ||
		f.ownerId != family.ownerId ||
		f.name != family.name ||
		f.haveJointAccount != family.haveJointAccount ||
		!f.createdAt.Equal(family.createdAt) ||
		!f.updatedAt.Equal(family.updatedAt) ||
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

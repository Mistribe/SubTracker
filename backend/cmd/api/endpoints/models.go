package endpoints

import (
	"errors"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/user"
)

type OwnerModel struct {
	Type     string  `json:"type" binding:"required"`
	FamilyId *string `json:"family_id,omitempty"`
	UserId   *string `json:"userId,omitempty"`
	Etag     string  `json:"etag" binding:"required"`
}

type editableOwnerModel struct {
	Type     string  `json:"type" binding:"required"`
	FamilyId *string `json:"family_id,omitempty"`
}

func (m editableOwnerModel) Owner(userId string) (user.Owner, error) {
	ownerType, err := user.ParseOwnerType(m.Type)
	if err != nil {
		return nil, err
	}

	switch ownerType {
	case user.PersonalOwner:
		return user.NewPersonalOwner(userId), nil
	case user.FamilyOwner:
		if m.FamilyId == nil {
			return nil, errors.New("missing family_id")
		}
		familyId, err2 := uuid.Parse(*m.FamilyId)
		if err2 != nil {
			return nil, err2
		}
		return user.NewFamilyOwner(familyId), nil
	}

	return nil, errors.New("unknown owner type")
}

func newOwnerModel(source user.Owner) OwnerModel {
	model := OwnerModel{
		Type: source.Type().String(),
		Etag: source.ETag(),
	}
	switch source.Type() {
	case user.PersonalOwner:
		userId := source.UserId()
		model.UserId = &userId
		break
	case user.FamilyOwner:
		familyId := source.FamilyId().String()
		model.FamilyId = &familyId
		break
	}

	return model
}

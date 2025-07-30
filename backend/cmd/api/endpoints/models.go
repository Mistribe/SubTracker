package endpoints

import (
	"errors"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
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

func (m editableOwnerModel) Owner(userId string) (auth.Owner, error) {
	ownerType, err := auth.ParseOwnerType(m.Type)
	if err != nil {
		return nil, err
	}

	switch ownerType {
	case auth.PersonalOwner:
		return auth.NewPersonalOwner(userId), nil
	case auth.FamilyOwner:
		if m.FamilyId == nil {
			return nil, errors.New("missing family_id")
		}
		familyId, err2 := uuid.Parse(*m.FamilyId)
		if err2 != nil {
			return nil, err2
		}
		return auth.NewFamilyOwner(familyId), nil
	}

	return nil, errors.New("unknown owner type")
}

func newOwnerModel(source auth.Owner) OwnerModel {
	model := OwnerModel{
		Type: source.Type().String(),
		Etag: source.ETag(),
	}
	switch source.Type() {
	case auth.PersonalOwner:
		userId := source.UserId()
		model.UserId = &userId
	case auth.FamilyOwner:
		familyId := source.FamilyId().String()
		model.FamilyId = &familyId
	}

	return model
}

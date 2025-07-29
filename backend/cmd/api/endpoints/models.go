package endpoints

import (
	"github.com/oleexo/subtracker/internal/domain/user"
)

type ownerModel struct {
	Type     string  `json:"type" binding:"required"`
	FamilyId *string `json:"family_id,omitempty"`
	UserId   *string `json:"userId,omitempty"`
	Etag     string  `json:"etag" binding:"required"`
}

func newOwnerModel(source user.Owner) ownerModel {
	model := ownerModel{
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

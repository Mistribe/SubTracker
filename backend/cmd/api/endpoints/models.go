package endpoints

import (
	"errors"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/currency"
)

// OwnerModel represents ownership information for resources
// @Description Owner object that can represent either personal or family ownership
type OwnerModel struct {
	// @Description Type of ownership (personal, family or system)
	Type string `json:"type" binding:"required" example:"personal" enums:"personal,family,system"`
	// @Description Family ID when an ownership type is family (required for family ownership)
	FamilyId *string `json:"family_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description UserProfile ID when an ownership type is personal (required for personal ownership)
	UserId *string `json:"userId,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
	// @Description Entity tag for optimistic concurrency control
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

type EditableOwnerModel struct {
	// @Description Type of ownership (personal, family or system)
	Type string `json:"type" binding:"required" example:"personal" enums:"personal,family,system"`
	// @Description Family ID when an ownership type is family (required for family ownership)
	FamilyId *string `json:"family_id,omitempty"`
}

func (m EditableOwnerModel) Owner(userId string) (auth.Owner, error) {
	ownerType, err := auth.ParseOwnerType(m.Type)
	if err != nil {
		return nil, err
	}

	switch ownerType {
	case auth.PersonalOwnerType:
		return auth.NewPersonalOwner(userId), nil
	case auth.FamilyOwnerType:
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
	case auth.PersonalOwnerType:
		userId := source.UserId()
		model.UserId = &userId
	case auth.FamilyOwnerType:
		familyId := source.FamilyId().String()
		model.FamilyId = &familyId
	}

	return model
}

type AmountModel struct {
	Value    float64 `json:"value" binding:"required" example:"100.00"`
	Currency string  `json:"currency" binding:"required" example:"USD"`
}

func newAmount(amount currency.Amount) AmountModel {
	return AmountModel{
		Value:    amount.Value(),
		Currency: amount.Currency().String(),
	}
}

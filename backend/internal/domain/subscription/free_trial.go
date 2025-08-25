package subscription

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type FreeTrial interface {
	entity.ETagEntity

	StartDate() time.Time
	EndDate() time.Time

	GetValidationErrors() validation.Errors
}

func NewFreeTrial(
	startDate time.Time,
	endDate time.Time) FreeTrial {
	return &freeTrial{
		startDate: startDate,
		endDate:   endDate,
	}
}

type freeTrial struct {
	startDate time.Time
	endDate   time.Time
}

func (f freeTrial) StartDate() time.Time {
	return f.startDate
}

func (f freeTrial) EndDate() time.Time {
	return f.endDate
}

func (f freeTrial) ETagFields() []interface{} {
	return []interface{}{
		f.endDate,
		f.startDate,
	}
}

func (f freeTrial) ETag() string {
	return entity.CalculateETag(f)
}

func (f freeTrial) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if f.startDate.IsZero() {
		errors = append(errors, validation.NewError("startDate", "StartDate is required"))
	}

	if f.endDate.IsZero() {
		errors = append(errors, validation.NewError("endDate", "EndDate is required"))
	}

	if !f.endDate.IsZero() && !f.startDate.IsZero() && f.endDate.Before(f.startDate) {
		errors = append(errors, validation.NewError("endDate", "EndDate must be after StartDate"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

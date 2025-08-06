package subscription

import (
	"time"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/validationx"
)

type FreeTrial interface {
	entity.ETagEntity

	StartDate() time.Time
	EndDate() time.Time

	GetValidationErrors() validationx.Errors
}

func NewFreeTrial(startDate time.Time,
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

func (f freeTrial) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if f.startDate.IsZero() {
		errors = append(errors, validationx.NewError("startDate", "StartDate is required"))
	}

	if f.endDate.IsZero() {
		errors = append(errors, validationx.NewError("endDate", "EndDate is required"))
	}

	if !f.endDate.IsZero() && !f.startDate.IsZero() && f.endDate.Before(f.startDate) {
		errors = append(errors, validationx.NewError("endDate", "EndDate must be after StartDate"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

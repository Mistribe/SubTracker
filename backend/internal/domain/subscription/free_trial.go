package subscription

import (
	"time"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type FreeTrial interface {
	entity.ETagEntity

	StartDate() time.Time
	EndDate() time.Time
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

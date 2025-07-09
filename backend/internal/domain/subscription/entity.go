package subscription

import (
	"time"

	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core/option"
)

type Subscription struct {
	id            uuid.UUID
	name          string
	payments      []Payment
	labels        []uuid.UUID
	familyMembers []uuid.UUID
	payer         option.Option[uuid.UUID]
	createdAt     time.Time
	updatedAt     time.Time
}

func (s Subscription) Id() uuid.UUID {
	return s.id
}

func (s Subscription) Name() string {
	return s.name
}

func (s Subscription) Payments() []Payment {
	return s.payments
}

func (s Subscription) Labels() []uuid.UUID {
	return s.labels
}

func (s Subscription) FamilyMembers() []uuid.UUID {
	return s.familyMembers
}

func (s Subscription) Payer() option.Option[uuid.UUID] {
	return s.payer
}

func (s Subscription) CreatedAt() time.Time {
	return s.createdAt
}

func (s Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

type Payment struct {
	id        uuid.UUID
	price     float64
	startDate time.Time
	endDate   option.Option[time.Time]
	months    int
	currency  string
	createdAt time.Time
	updatedAt time.Time
}

func (p Payment) Id() uuid.UUID {
	return p.id
}

func (p Payment) Price() float64 {
	return p.price
}

func (p Payment) StartDate() time.Time {
	return p.startDate
}

func (p Payment) EndDate() option.Option[time.Time] {
	return p.endDate
}

func (p Payment) Months() int {
	return p.months
}

func (p Payment) Currency() string {
	return p.currency
}

func (p Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

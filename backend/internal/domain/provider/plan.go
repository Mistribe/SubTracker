package provider

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"github.com/oleexo/subtracker/pkg/validationx"
)

func planUniqueComparer(p1, p2 Plan) bool {
	return p1.Id() == p2.Id()
}

func planComparer(p1, p2 Plan) bool {
	return p1.Equal(p2)
}

type Plan interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	SetName(name string)
	Description() *string
	SetDescription(description *string)
	Prices() *slicesx.Tracked[Price]
	SetPrices(prices []Price)
	Equal(other Plan) bool
	GetValidationErrors() validationx.Errors
	ContainsPrice(priceId uuid.UUID) bool
	AddPrice(price Price) bool
	GetPriceById(priceId uuid.UUID) Price
	RemovePriceById(priceId uuid.UUID) bool
}

type plan struct {
	*entity.Base

	name        string
	description *string
	prices      *slicesx.Tracked[Price]
}

func NewPlan(
	id uuid.UUID,
	name string,
	description *string,
	prices []Price,
	createdAt time.Time,
	updatedAt time.Time,
) Plan {
	return &plan{
		Base:        entity.NewBase(id, createdAt, updatedAt, true, false),
		name:        name,
		description: description,
		prices:      slicesx.NewTracked(prices, priceUniqueComparer, priceComparer),
	}
}

func (p *plan) RemovePriceById(priceId uuid.UUID) bool {
	pr := p.GetPriceById(priceId)
	if pr == nil {
		return false
	}

	return p.prices.Remove(pr)
}

func (p *plan) AddPrice(price Price) bool {
	if p.prices.Add(price) {
		p.SetAsDirty()
		return true
	}
	return false
}

func (p *plan) GetPriceById(priceId uuid.UUID) Price {
	for pr := range p.prices.It() {
		if pr.Id() == priceId {
			return pr
		}
	}

	return nil
}

func (p *plan) ContainsPrice(priceId uuid.UUID) bool {
	for pr := range p.prices.It() {
		if pr.Id() == priceId {
			return true
		}
	}

	return false
}
func (p *plan) Name() string {
	return p.name
}

func (p *plan) SetName(name string) {
	p.name = name
	p.SetAsDirty()
}

func (p *plan) Description() *string {
	return p.description
}

func (p *plan) SetDescription(description *string) {
	p.description = description
	p.SetAsDirty()
}

func (p *plan) Prices() *slicesx.Tracked[Price] {
	return p.prices
}

func (p *plan) SetPrices(prices []Price) {
	p.prices = slicesx.NewTracked(prices, priceUniqueComparer, priceComparer)
	p.SetAsDirty()
}

func (p *plan) ETagFields() []interface{} {
	fields := []interface{}{
		p.name,
		p.description,
	}

	for pr := range p.prices.It() {
		fields = append(fields, pr.ETagFields()...)
	}

	return fields
}

func (p *plan) ETag() string {
	return entity.CalculateETag(p, p.Base)
}

func (p *plan) Equal(other Plan) bool {
	if other == nil {
		return false
	}

	return p.ETag() == other.ETag()
}

func (p *plan) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if strings.TrimSpace(p.name) == "" {
		errors = append(errors, validationx.NewError("name", "name is required and cannot be empty"))
	}

	if p.prices == nil || len(p.prices.Values()) == 0 {
		errors = append(errors, validationx.NewError("prices", "at least one price is required"))
	}

	for _, pr := range p.prices.Values() {
		if err := pr.GetValidationErrors(); err != nil {
			errors = append(errors, err...)
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

package provider

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"github.com/oleexo/subtracker/pkg/validationx"
	"github.com/oleexo/subtracker/pkg/x"
)

type Provider interface {
	entity.Entity
	entity.ETagEntity

	Name() string
	SetName(name string)
	Description() *string
	SetDescription(description *string)
	IconUrl() *string
	SetIconUrl(iconUrl *string)
	Url() *string
	SetUrl(url *string)
	PricingPageUrl() *string
	SetPricingPageUrl(pricingUrl *string)
	Labels() *slicesx.Tracked[uuid.UUID]
	SetLabels(labels []uuid.UUID)
	IsCustom() bool
	SetOwner(owner auth.Owner)
	Owner() auth.Owner
	Plans() *slicesx.Tracked[Plan]
	SetPlans(plans []Plan)
	Equal(other Provider) bool
	GetValidationErrors() validationx.Errors
	ContainsPlan(planId uuid.UUID) bool
	AddPlan(plan Plan) bool
	GetPlanById(planId uuid.UUID) Plan
	RemovePlanById(planId uuid.UUID) bool
}

type provider struct {
	*entity.Base

	name           string
	description    *string
	iconUrl        *string
	url            *string
	pricingPageUrl *string
	labels         *slicesx.Tracked[uuid.UUID]
	plans          *slicesx.Tracked[Plan]
	owner          auth.Owner
}

func NewProvider(
	id uuid.UUID,
	name string,
	description *string,
	iconUrl *string,
	url *string,
	pricingPageUrl *string,
	labels []uuid.UUID,
	plans []Plan,
	owner auth.Owner,
	createdAt time.Time,
	updatedAt time.Time) Provider {
	return &provider{
		Base:           entity.NewBase(id, createdAt, updatedAt, true, false),
		name:           name,
		description:    description,
		iconUrl:        iconUrl,
		url:            url,
		pricingPageUrl: pricingPageUrl,
		labels:         slicesx.NewTracked(labels, x.UuidUniqueComparer, x.UuidComparer),
		plans:          slicesx.NewTracked(plans, planUniqueComparer, planComparer),
		owner:          owner,
	}
}

func (p *provider) RemovePlanById(planId uuid.UUID) bool {
	pl := p.GetPlanById(planId)
	if pl == nil {
		return false
	}

	return p.plans.Remove(pl)
}

func (p *provider) GetPlanById(planId uuid.UUID) Plan {
	for pl := range p.plans.It() {
		if pl.Id() == planId {
			return pl
		}
	}

	return nil
}

func (p *provider) ContainsPlan(planId uuid.UUID) bool {
	for pl := range p.plans.It() {
		if pl.Id() == planId {
			return true
		}
	}

	return false
}

func (p *provider) AddPlan(plan Plan) bool {
	if p.plans.Add(plan) {
		p.SetAsDirty()
		return true
	}
	return false
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) SetName(name string) {
	p.name = name
	p.SetAsDirty()
}

func (p *provider) Description() *string {
	return p.description
}

func (p *provider) SetDescription(description *string) {
	p.description = description
	p.SetAsDirty()
}

func (p *provider) IconUrl() *string {
	return p.iconUrl
}

func (p *provider) SetIconUrl(iconUrl *string) {
	p.iconUrl = iconUrl
	p.SetAsDirty()
}

func (p *provider) Url() *string {
	return p.url
}

func (p *provider) SetUrl(url *string) {
	p.url = url
	p.SetAsDirty()
}

func (p *provider) PricingPageUrl() *string {
	return p.pricingPageUrl
}

func (p *provider) SetPricingPageUrl(pricingUrl *string) {
	p.pricingPageUrl = pricingUrl
	p.SetAsDirty()
}

func (p *provider) Labels() *slicesx.Tracked[uuid.UUID] {
	return p.labels
}

func (p *provider) SetLabels(labels []uuid.UUID) {
	p.labels = slicesx.NewTracked(labels, x.UuidUniqueComparer, x.UuidComparer)
	p.SetAsDirty()
}

func (p *provider) IsCustom() bool {
	return p.owner != nil
}

func (p *provider) SetOwner(owner auth.Owner) {
	p.owner = owner
	p.SetAsDirty()
}

func (p *provider) Owner() auth.Owner {
	return p.owner
}

func (p *provider) Plans() *slicesx.Tracked[Plan] {
	return p.plans
}

func (p *provider) SetPlans(plans []Plan) {
	p.plans = slicesx.NewTracked(plans, planUniqueComparer, planComparer)
	p.SetAsDirty()
}

func (p *provider) ETagFields() []interface{} {
	fields := []interface{}{
		p.name,
		p.description,
		p.iconUrl,
		p.url,
		p.pricingPageUrl,
		p.owner.ETag(),
	}

	for pl := range p.plans.It() {
		fields = append(fields, pl.ETagFields()...)
	}

	for lbl := range p.labels.It() {
		fields = append(fields, lbl.String())
	}

	return fields
}

func (p *provider) ETag() string {
	return entity.CalculateETag(p, p.Base)
}

func (p *provider) Equal(other Provider) bool {
	if other == nil {
		return false
	}

	return p.ETag() == other.ETag()
}

func (p *provider) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if strings.TrimSpace(p.name) == "" {
		errors = append(errors, validationx.NewError(
			"name",
			"name is required and cannot be empty",
		))
	}

	if p.plans != nil {
		for _, pl := range p.plans.Values() {
			if planErr := pl.GetValidationErrors(); planErr != nil {
				errors = append(errors, planErr...)
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

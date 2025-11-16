package provider

import (
	"strings"
	"time"

	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type Provider interface {
	entity.Entity[types.ProviderID]
	entity.ETagEntity

	Name() string
	Key() string
	SetName(name string)
	Description() *string
	SetDescription(description *string)
	IconUrl() *string
	SetIconUrl(iconUrl *string)
	Url() *string
	SetUrl(url *string)
	PricingPageUrl() *string
	SetPricingPageUrl(pricingUrl *string)
	Labels() *slicesx.Tracked[types.LabelID]
	SetLabels(labels []types.LabelID)
	IsCustom() bool
	SetOwner(owner types.Owner)
	Owner() types.Owner
	Equal(other Provider) bool
	GetValidationErrors() validation.Errors
}

type provider struct {
	*entity.Base[types.ProviderID]

	name           string
	key            string
	description    *string
	iconUrl        *string
	url            *string
	pricingPageUrl *string
	labels         *slicesx.Tracked[types.LabelID]
	owner          types.Owner
}

func NewProvider(
	id types.ProviderID,
	name string,
	description *string,
	iconUrl *string,
	url *string,
	pricingPageUrl *string,
	labels []types.LabelID,
	owner types.Owner,
	createdAt time.Time,
	updatedAt time.Time) Provider {

	return &provider{
		Base:           entity.NewBase(id, createdAt, updatedAt, true, false),
		name:           strings.TrimSpace(name),
		key:            generateKey(name, owner),
		description:    description,
		iconUrl:        iconUrl,
		url:            url,
		pricingPageUrl: pricingPageUrl,
		labels:         slicesx.NewTracked[types.LabelID](labels, types.LabelIDComparer, types.LabelIDComparer),
		owner:          owner,
	}
}

func (p *provider) Key() string {
	return p.key
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) SetName(name string) {
	if name == p.name {
		return
	}
	p.name = name
	p.key = generateKey(name, p.owner)
	p.SetAsDirty()
}

func (p *provider) Description() *string {
	return p.description
}

func (p *provider) SetDescription(description *string) {
	if description == p.description {
		return
	}
	p.description = description
	p.SetAsDirty()
}

func (p *provider) IconUrl() *string {
	return p.iconUrl
}

func (p *provider) SetIconUrl(iconUrl *string) {
	if iconUrl == p.iconUrl {
		return
	}
	p.iconUrl = iconUrl
	p.SetAsDirty()
}

func (p *provider) Url() *string {
	return p.url
}

func (p *provider) SetUrl(url *string) {
	if url == p.url {
		return
	}
	p.url = url
	p.SetAsDirty()
}

func (p *provider) PricingPageUrl() *string {
	return p.pricingPageUrl
}

func (p *provider) SetPricingPageUrl(pricingUrl *string) {
	if p.pricingPageUrl == pricingUrl {
		return
	}
	p.pricingPageUrl = pricingUrl
	p.SetAsDirty()
}

func (p *provider) Labels() *slicesx.Tracked[types.LabelID] {
	return p.labels
}

func (p *provider) SetLabels(labels []types.LabelID) {
	p.labels.Set(labels)
	p.SetAsDirty()
}

func (p *provider) IsCustom() bool {
	return p.owner != nil
}

func (p *provider) SetOwner(owner types.Owner) {
	p.owner = owner
	p.key = generateKey(p.name, owner)
	p.SetAsDirty()
}

func (p *provider) Owner() types.Owner {
	return p.owner
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

	for lbl := range p.labels.It() {
		fields = append(fields, lbl.String())
	}

	return fields
}

func (p *provider) ETag() string {
	return entity.CalculateETag(p)
}

func (p *provider) Equal(other Provider) bool {
	if other == nil {
		return false
	}

	return p.ETag() == other.ETag()
}

func (p *provider) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if strings.TrimSpace(p.name) == "" {
		errors = append(errors, validation.NewError(
			"name",
			"name is required and cannot be empty",
		))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

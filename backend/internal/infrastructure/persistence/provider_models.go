package persistence

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type providerPriceSqlModel struct {
	baseSqlModel

	StartDate time.Time
	EndDate   *time.Time
	Currency  string
	Amount    float64
	PlanId    uuid.UUID
}

func (p providerPriceSqlModel) TableName() string {
	return "provider_prices"
}

func newProviderPriceSqlModel(planId uuid.UUID, source provider.Price) providerPriceSqlModel {
	return providerPriceSqlModel{
		baseSqlModel: newBaseSqlModel(source),
		StartDate:    source.StartDate(),
		EndDate:      source.EndDate(),
		Currency:     source.Currency().String(),
		Amount:       source.Amount(),
		PlanId:       planId,
	}
}

func newProviderPrice(model providerPriceSqlModel) provider.Price {
	return provider.NewPrice(
		model.Id,
		model.StartDate,
		model.EndDate,
		currency.MustParseISO(model.Currency),
		model.Amount,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

type providerPlanSqlModel struct {
	baseSqlModel

	Name        *string                 `gorm:"type:varchar(100)"`
	Description *string                 `gorm:"type:varchar(255)"`
	Prices      []providerPriceSqlModel `gorm:"foreignKey:PlanId;references:Id;constraint:OnDelete:CASCADE"`
	ProviderId  uuid.UUID               `gorm:"type:uuid;not null"`
}

func (p providerPlanSqlModel) TableName() string {
	return "provider_plans"
}

func newProviderPlanSqlModel(providerId uuid.UUID, source provider.Plan) providerPlanSqlModel {
	return providerPlanSqlModel{
		baseSqlModel: newBaseSqlModel(source),
		Name:         source.Name(),
		Description:  source.Description(),
		ProviderId:   providerId,
	}
}

func newProviderPlan(model providerPlanSqlModel) provider.Plan {
	var prices []provider.Price
	if model.Prices != nil && len(model.Prices) > 0 {
		prices = make([]provider.Price, len(model.Prices))
		for i, price := range model.Prices {
			prices[i] = newProviderPrice(price)
		}
	}
	return provider.NewPlan(
		model.Id,
		model.Name,
		model.Description,
		prices,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

type providerLabelSqlModel struct {
	LabelId    uuid.UUID        `gorm:"primaryKey;type:uuid;not null"`
	Label      labelSqlModel    `gorm:"foreignKey:LabelId;references:Id"`
	ProviderId uuid.UUID        `gorm:"primaryKey;type:uuid;not null"`
	Provider   providerSqlModel `gorm:"foreignKey:ProviderId;references:Id"`
}

func (p providerLabelSqlModel) TableName() string {
	return "provider_labels"
}

func newProviderLabelSqlModel(providerId uuid.UUID, labelId uuid.UUID) providerLabelSqlModel {
	return providerLabelSqlModel{
		LabelId:    labelId,
		ProviderId: providerId,
	}
}

type providerSqlModel struct {
	baseSqlModel
	baseOwnerSqlModel

	Name           string                  `gorm:"type:varchar(100);not null"`
	Description    *string                 `gorm:"type:varchar(255)"`
	IconUrl        *string                 `gorm:"type:varchar(255)"`
	Url            *string                 `gorm:"type:varchar(255)"`
	PricingPageUrl *string                 `gorm:"type:varchar(255)"`
	Labels         []providerLabelSqlModel `gorm:"foreignKey:ProviderId;references:Id;constraint:OnDelete:CASCADE"`
	Plans          []providerPlanSqlModel  `gorm:"foreignKey:ProviderId;references:Id;constraint:OnDelete:CASCADE"`
}

func (p providerSqlModel) TableName() string {
	return "providers"
}

func newProviderSqlModel(source provider.Provider) providerSqlModel {
	model := providerSqlModel{
		baseSqlModel:   newBaseSqlModel(source),
		Name:           source.Name(),
		Description:    source.Description(),
		IconUrl:        source.IconUrl(),
		Url:            source.Url(),
		PricingPageUrl: source.PricingPageUrl(),
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case user.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
		break
	case user.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = &userId
		break
	}
	return model
}

func newProvider(model providerSqlModel) provider.Provider {
	ownerType, err := user.ParseOwnerType(model.OwnerType)
	if err != nil {
		panic(err)
	}
	owner := user.NewOwner(ownerType, model.OwnerFamilyId, model.OwnerUserId)
	var labels []uuid.UUID
	if model.Labels != nil && len(model.Labels) > 0 {
		labels = make([]uuid.UUID, len(model.Labels))
		for i, label := range model.Labels {
			labels[i] = label.LabelId
		}
	}

	var plans []provider.Plan
	if model.Plans != nil && len(model.Plans) > 0 {
		plans = make([]provider.Plan, len(model.Plans))
		for i, plan := range model.Plans {
			plans[i] = newProviderPlan(plan)
		}
	}

	return provider.NewProvider(
		model.Id,
		model.Name,
		model.Description,
		model.IconUrl,
		model.Url,
		model.PricingPageUrl,
		labels,
		plans,
		owner,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

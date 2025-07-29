package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type providerPriceSqlModel struct {
	BaseSqlModel `gorm:"embedded"`

	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	Currency  string    `gorm:"type:varchar(3);not null"`
	Amount    float64   `gorm:"not null"`
	PlanId    uuid.UUID `gorm:"type:uuid;not null"`
}

func (p providerPriceSqlModel) TableName() string {
	return "provider_prices"
}

func newProviderPriceSqlModel(planId uuid.UUID, source provider.Price) providerPriceSqlModel {
	return providerPriceSqlModel{
		BaseSqlModel: newBaseSqlModel(source),
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
	BaseSqlModel `gorm:"embedded"`

	Name        sql.NullString          `gorm:"type:varchar(100)"`
	Description sql.NullString          `gorm:"type:varchar(255)"`
	Prices      []providerPriceSqlModel `gorm:"foreignKey:PlanId;references:Id;constraint:OnDelete:CASCADE"`
	ProviderId  uuid.UUID               `gorm:"type:uuid;not null"`
}

func (p providerPlanSqlModel) TableName() string {
	return "provider_plans"
}

func newProviderPlanSqlModel(providerId uuid.UUID, source provider.Plan) providerPlanSqlModel {
	model := providerPlanSqlModel{
		BaseSqlModel: newBaseSqlModel(source),
		ProviderId:   providerId,
	}

	if source.Name() != nil {
		model.Name = sql.NullString{
			String: *source.Name(),
			Valid:  true,
		}
	} else {
		model.Name = sql.NullString{
			Valid: false,
		}
	}

	if source.Description() != nil {
		model.Description = sql.NullString{
			String: *source.Description(),
			Valid:  true,
		}
	} else {
		model.Description = sql.NullString{
			Valid: false,
		}
	}

	return model
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
		sqlNullToString(model.Name),
		sqlNullToString(model.Description),
		prices,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

type providerLabelSqlModel struct {
	LabelId    uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	ProviderId uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
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

type ProviderSqlModel struct {
	BaseSqlModel      `gorm:"embedded"`
	BaseOwnerSqlModel `gorm:"embedded"`

	Name           string                  `gorm:"type:varchar(100);not null"`
	Description    sql.NullString          `gorm:"type:varchar(255)"`
	IconUrl        sql.NullString          `gorm:"type:varchar(255)"`
	Url            sql.NullString          `gorm:"type:varchar(255)"`
	PricingPageUrl sql.NullString          `gorm:"type:varchar(255)"`
	Labels         []providerLabelSqlModel `gorm:"many2many:provider_labels;joinForeignKey:ProviderId;joinReferences:LabelId"`
	Plans          []providerPlanSqlModel  `gorm:"foreignKey:ProviderId;references:Id;constraint:OnDelete:CASCADE"`
}

func (p ProviderSqlModel) TableName() string {
	return "providers"
}

func newProviderSqlModel(source provider.Provider) ProviderSqlModel {
	model := ProviderSqlModel{
		BaseSqlModel:   newBaseSqlModel(source),
		Name:           source.Name(),
		Description:    stringToSqlNull(source.Description()),
		IconUrl:        stringToSqlNull(source.IconUrl()),
		Url:            stringToSqlNull(source.Url()),
		PricingPageUrl: stringToSqlNull(source.PricingPageUrl()),
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case user.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
		break
	case user.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = stringToSqlNull(&userId)
		break
	}
	return model
}

func newProvider(model ProviderSqlModel) provider.Provider {
	ownerType, err := user.ParseOwnerType(model.OwnerType)
	if err != nil {
		panic(err)
	}
	owner := user.NewOwner(ownerType, model.OwnerFamilyId, sqlNullToString(model.OwnerUserId))
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
		sqlNullToString(model.Description),
		sqlNullToString(model.IconUrl),
		sqlNullToString(model.Url),
		sqlNullToString(model.PricingPageUrl),
		labels,
		plans,
		owner,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

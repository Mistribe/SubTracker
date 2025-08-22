package endpoints

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"github.com/oleexo/subtracker/pkg/x"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func (s SubscriptionEndpointGroup) Prefix() string {
	return "/subscriptions"
}

func (s SubscriptionEndpointGroup) Routes() []ginfx.Route {
	return s.routes
}

func (s SubscriptionEndpointGroup) Middlewares() []gin.HandlerFunc {
	return s.middlewares
}

func NewSubscriptionEndpointGroup(
	getEndpoint *SubscriptionGetEndpoint,
	getAllEndpoint *SubscriptionGetAllEndpoint,
	createEndpoint *SubscriptionCreateEndpoint,
	updateEndpoint *SubscriptionUpdateEndpoint,
	deleteEndpoint *SubscriptionDeleteEndpoint,
	patchEndpoint *SubscriptionPatchEndpoint,
	summaryEndpoint *SubscriptionSummaryEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
			summaryEndpoint,
			getEndpoint,
			getAllEndpoint,
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			patchEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

type SubscriptionFreeTrialModel struct {
	StartDate time.Time `json:"start_date" bindind:"required" format:"date-time"`
	EndDate   time.Time `json:"end_date" bindind:"required" format:"date-time"`
}

func newSubscriptionFreeTrialModel(source subscription.FreeTrial) *SubscriptionFreeTrialModel {
	if source == nil {
		return nil
	}

	return &SubscriptionFreeTrialModel{
		StartDate: source.StartDate(),
		EndDate:   source.EndDate(),
	}
}

func newSubscriptionFreeTrial(model *SubscriptionFreeTrialModel) subscription.FreeTrial {
	if model == nil {
		return nil
	}

	return subscription.NewFreeTrial(model.StartDate, model.EndDate)
}

func newSubscriptionCustomPrice(model *AmountModel) (subscription.CustomPrice, error) {
	if model == nil {
		return nil, nil
	}

	cry, err := currency.ParseISO(model.Currency)
	if err != nil {
		return nil, err
	}

	return subscription.NewCustomPrice(model.Value, cry), nil

}

// SubscriptionPayerModel represents who pays for a subscription within a family
// @Description Subscription payer object defining which family member is responsible for payment
type SubscriptionPayerModel struct {
	// @Description Type of payer (family or family member)
	Type string `json:"type" binding:"required" example:"family_member" enums:"family,family_member"`
	// @Description ID of the family associated with this payer
	FamilyId string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description ID of the specific family member who pays (required when type is family_member)
	MemberId *string `json:"memberId,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

// EditableSubscriptionPayerModel represents editable payer information for subscription updates
// @Description Subscription payer object used for updating who pays for a subscription
type EditableSubscriptionPayerModel struct {
	// @Description Type of payer (family or family member)
	Type string `json:"type" binding:"required" example:"family_member" enums:"family,family_member"`
	// @Description ID of the family associated with this payer
	FamilyId string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description ID of the specific family member who pays (required when type is family_member)
	MemberId *string `json:"memberId,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// SubscriptionModel represents an active subscription to a service provider's plan
// @Description Subscription object containing all information about an active subscription including billing and usage details
type SubscriptionModel struct {
	// @Description Unique identifier for the subscription (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Optional custom name for easy identification of the subscription
	FriendlyName *string `json:"friendly_name,omitempty" example:"Netflix Family Account" maxLength:"255"`
	// @Description Number of free trial days remaining (null if no trial or trial expired)
	FreeTrial *SubscriptionFreeTrialModel `json:"free_trial,omitempty"`
	// @Description ID of the service provider offering this subscription
	ProviderId string `json:"provider_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174002"`
	// @Description ID of the specific plan being subscribed to
	PlanId *string `json:"plan_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174003"`
	// @Description ID of the pricing tier for this subscription
	PriceId     *string      `json:"price_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174004"`
	CustomPrice *AmountModel `json:"custom_price,omitempty"`
	// @Description Ownership information specifying whether this subscription belongs to a user or family
	Owner OwnerModel `json:"owner" binding:"required"`
	// @Description List of family member IDs who use this service (for shared subscriptions)
	ServiceUsers []string `json:"service_users,omitempty" example:"123e4567-e89b-12d3-a456-426614174005,123e4567-e89b-12d3-a456-426614174006"`
	// @Description ISO 8601 timestamp when the subscription becomes active
	StartDate time.Time `json:"start_date" binding:"required" format:"date-time" example:"2023-01-01T00:00:00Z"`
	// @Description ISO 8601 timestamp when the subscription expires (null for ongoing subscriptions)
	EndDate *time.Time `json:"end_date,omitempty" format:"date-time" example:"2024-01-01T00:00:00Z"`
	// @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
	Recurrency string `json:"recurrency" binding:"required" example:"monthly" enums:"unknown,one_time,monthly,quarterly,half_yearly,yearly,custom"`
	// @Description CustomRecurrency recurrency interval in days (required when recurrency is custom)
	CustomRecurrency *int32 `json:"custom_recurrency,omitempty" example:"90" minimum:"1" maximum:"3650"`
	// @Description Information about who pays for this subscription within the family
	Payer *SubscriptionPayerModel `json:"payer,omitempty"`
	// @Description ISO 8601 timestamp when the subscription was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp when the subscription was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
	// @Description Indicates whether the subscription is currently active or not
	IsActive bool `json:"is_active" binding:"required" example:"true"`
	// @Description List of labels associated with this subscription
	LabelRefs []LabelRefModel `json:"label_refs,omitempty"`
}

func newSubscriptionPayerModel(source subscription.Payer) SubscriptionPayerModel {
	model := SubscriptionPayerModel{
		Type:     source.Type().String(),
		FamilyId: source.FamilyId().String(),
		Etag:     source.ETag(),
	}

	if source.Type() == subscription.FamilyMemberPayer {
		memberId := source.MemberId().String()
		model.MemberId = &memberId
	}
	return model
}

type LabelRefModel struct {
	LabelId string `json:"label_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Source  string `json:"source" binding:"required" example:"subscription" enums:"subscription,provider"`
}

func newSubscriptionLabelRef(ref subscription.LabelRef) LabelRefModel {
	return LabelRefModel{
		LabelId: ref.LabelId.String(),
		Source:  ref.Source.String(),
	}
}

func newSubscriptionModel(source subscription.Subscription) SubscriptionModel {
	var payerModel SubscriptionPayerModel
	if source.Payer() != nil {
		payerModel = newSubscriptionPayerModel(source.Payer())
	}
	serviceUsers := slicesx.Select(source.ServiceUsers().Values(), func(in uuid.UUID) string {
		return in.String()
	})
	labelRefs := slicesx.Select(source.Labels().Values(), newSubscriptionLabelRef)
	model := SubscriptionModel{
		Id:               source.Id().String(),
		FriendlyName:     source.FriendlyName(),
		FreeTrial:        newSubscriptionFreeTrialModel(source.FreeTrial()),
		ProviderId:       source.ProviderId().String(),
		ServiceUsers:     serviceUsers,
		LabelRefs:        labelRefs,
		Owner:            newOwnerModel(source.Owner()),
		StartDate:        source.StartDate(),
		EndDate:          source.EndDate(),
		Recurrency:       source.Recurrency().String(),
		CustomRecurrency: source.CustomRecurrency(),
		Payer:            &payerModel,
		IsActive:         source.IsActive(),
		CreatedAt:        source.CreatedAt(),
		UpdatedAt:        source.UpdatedAt(),
		Etag:             source.ETag(),
		CustomPrice:      x.P(newAmount(source.CustomPrice())),
	}

	if source.PlanId() != nil {
		planId := source.PlanId().String()
		model.PlanId = &planId
	}

	if source.PriceId() != nil {
		priceId := source.PriceId().String()
		model.PriceId = &priceId
	}

	return model
}

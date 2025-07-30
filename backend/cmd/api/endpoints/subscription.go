package endpoints

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/pkg/slicesx"

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
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
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

type subscriptionPayerModel struct {
	Type     string  `json:"type" binding:"required"`
	FamilyId string  `json:"family_id" binding:"required"`
	MemberId *string `json:"memberId,omitempty"`
	Etag     string  `json:"etag" binding:"required"`
}

type editableSubscriptionPayerModel struct {
	Type     string  `json:"type" binding:"required"`
	FamilyId string  `json:"family_id" binding:"required"`
	MemberId *string `json:"memberId,omitempty"`
}

type subscriptionModel struct {
	Id                string                  `json:"id" binding:"required"`
	FriendlyName      *string                 `json:"friendly_name,omitempty"`
	FreeTrialDays     *uint                   `json:"free_trial_days,omitempty"`
	ServiceProviderId string                  `json:"service_provider_id" binding:"required"`
	PlanId            string                  `json:"plan_id" binding:"required"`
	PriceId      string     `json:"price_id" binding:"required"`
	Owner        OwnerModel `json:"owner" binding:"required"`
	ServiceUsers []string   `json:"service_users,omitempty"`
	StartDate         time.Time               `json:"start_date" binding:"required" format:"date-time"`
	EndDate           *time.Time              `json:"end_date,omitempty" format:"date-time"`
	Recurrency        string                  `json:"recurrency" binding:"required"`
	CustomRecurrency  *uint                   `json:"custom_recurrency,omitempty"`
	Payer             *subscriptionPayerModel `json:"payer,omitempty"`
	CreatedAt         time.Time               `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt         time.Time               `json:"updated_at" binding:"required" format:"date-time"`
	Etag              string                  `json:"etag" binding:"required"`
}

func newSubscriptionPayerModel(source subscription.Payer) subscriptionPayerModel {
	model := subscriptionPayerModel{
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

func newSubscriptionModel(source subscription.Subscription) subscriptionModel {
	var payerModel subscriptionPayerModel
	if source.Payer() != nil {
		payerModel = newSubscriptionPayerModel(source.Payer())
	}
	serviceUsers := slicesx.Map(source.ServiceUsers().Values(), func(in uuid.UUID) string {
		return in.String()
	})
	return subscriptionModel{
		Id:                source.Id().String(),
		FriendlyName:      source.FriendlyName(),
		FreeTrialDays:     source.FreeTrialDays(),
		ServiceProviderId: source.ServiceProviderId().String(),
		PlanId:            source.PlanId().String(),
		PriceId:           source.PriceId().String(),
		ServiceUsers:      serviceUsers,
		Owner:             newOwnerModel(source.Owner()),
		StartDate:         source.StartDate(),
		EndDate:           source.EndDate(),
		Recurrency:        source.Recurrency().String(),
		CustomRecurrency:  source.CustomRecurrency(),
		Payer:             &payerModel,
		CreatedAt:         source.CreatedAt(),
		UpdatedAt:         source.UpdatedAt(),
		Etag:              source.ETag(),
	}
}

package endpoints

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"

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
	createPaymentEndpoint *SubscriptionPaymentCreateEndpoint,
	updatePaymentEndpoint *SubscriptionPaymentUpdateEndpoint,
	deletePaymentEndpoint *SubscriptionPaymentDeleteEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
			getEndpoint,
			getAllEndpoint,
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			patchEndpoint,
			createPaymentEndpoint,
			updatePaymentEndpoint,
			deletePaymentEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

type paymentModel struct {
	Id        string     `json:"id" binding:"required"`
	Price     float64    `json:"price" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Months    int        `json:"months" binding:"required"`
	Currency  string     `json:"currency" binding:"required"`
	CreatedAt time.Time  `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt time.Time  `json:"updated_at" binding:"required" format:"date-time"`
	Etag      string     `json:"etag" binding:"required"`
}

type subscriptionModel struct {
	Id                  string         `json:"id" binding:"required"`
	Name                string         `json:"name" binding:"required"`
	Payments            []paymentModel `json:"payments" binding:"required"`
	Labels              []string       `json:"labels" binding:"required"`
	FamilyMembers       []string       `json:"family_members" binding:"required"`
	PayerId             *string        `json:"payer_id_id,omitempty"`
	FamilyId            *string        `json:"family_id,omitempty"`
	PayedByJointAccount bool           `json:"payed_by_joint_account" binding:"required"`
	CreatedAt           time.Time      `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt           time.Time      `json:"updated_at" binding:"required" format:"date-time"`
	Etag                string         `json:"etag" binding:"required"`
}

func newPaymentModel(source subscription.Payment) paymentModel {
	return paymentModel{
		Id:        source.Id().String(),
		Price:     source.Price(),
		StartDate: source.StartDate(),
		EndDate:   source.EndDate().Value(),
		Months:    source.Months(),
		Currency:  source.Currency().String(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}

func newSubscriptionModel(source subscription.Subscription) subscriptionModel {
	return subscriptionModel{
		Id:                  source.Id().String(),
		Name:                source.Name(),
		Payments:            ext.Map(source.Payments().Values(), newPaymentModel),
		Labels:              ext.Map(source.Labels().Values(), ext.UuidToString),
		FamilyMembers:       ext.Map(source.FamilyMembers().Values(), ext.UuidToString),
		PayerId:             option.Map[uuid.UUID, string](source.Payer(), ext.UuidToString).Value(),
		FamilyId:            option.Map[uuid.UUID, string](source.FamilyId(), ext.UuidToString).Value(),
		PayedByJointAccount: source.PayedByJointAccount(),
		CreatedAt:           source.CreatedAt(),
		UpdatedAt:           source.UpdatedAt(),
		Etag:                source.ETag(),
	}
}

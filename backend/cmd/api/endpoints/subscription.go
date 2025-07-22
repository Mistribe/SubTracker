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
	Id        string     `json:"id"`
	Price     float64    `json:"price"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Months    int        `json:"months"`
	Currency  string     `json:"currency"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Etag      string     `json:"etag"`
}

type subscriptionModel struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Payments      []paymentModel `json:"payments"`
	Labels        []string       `json:"labels"`
	FamilyMembers []string       `json:"family_members"`
	PayerId       *string        `json:"payer_id_id,omitempty"`
	FamilyId      *string        `json:"family_id,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Etag          string         `json:"etag"`
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
		Id:            source.Id().String(),
		Name:          source.Name(),
		Payments:      ext.Map(source.Payments().Values(), newPaymentModel),
		Labels:        ext.Map(source.Labels().Values(), ext.UuidToString),
		FamilyMembers: ext.Map(source.FamilyMembers().Values(), ext.UuidToString),
		PayerId:       option.Map[uuid.UUID, string](source.Payer(), ext.UuidToString).Value(),
		FamilyId:      option.Map[uuid.UUID, string](source.FamilyId(), ext.UuidToString).Value(),
		CreatedAt:     source.CreatedAt(),
		UpdatedAt:     source.UpdatedAt(),
		Etag:          source.ETag(),
	}
}

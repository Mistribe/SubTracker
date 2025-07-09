package endpoints

import (
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/pkg/ext"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/api/ginfx"
)

type SubscriptionEndpointGroup struct {
	routes []ginfx.Route
}

func (s SubscriptionEndpointGroup) Prefix() string {
	return "/subscription"
}

func (s SubscriptionEndpointGroup) Routes() []ginfx.Route {
	return s.routes
}

func (s SubscriptionEndpointGroup) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionEndpointGroup(getEndpoint *SubscriptionGetEndpoint,
	getAllEndpoint *SubscriptionGetAllEndpoint) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
			getEndpoint,
			getAllEndpoint,
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
}

type subscriptionModel struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Payments      []paymentModel `json:"payments"`
	Labels        []string       `json:"labels"`
	FamilyMembers []string       `json:"family_members"`
	Payer         *string        `json:"payer,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func newPaymentModel(source subscription.Payment) paymentModel {
	return paymentModel{
		Id:        source.Id().String(),
		Price:     source.Price(),
		StartDate: source.StartDate(),
		EndDate:   source.EndDate().Value(),
		Months:    source.Months(),
		Currency:  source.Currency(),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
	}
}

func newSubscriptionModel(source subscription.Subscription) subscriptionModel {
	return subscriptionModel{
		Id:            source.Id().String(),
		Name:          source.Name(),
		Payments:      ext.Map(source.Payments(), newPaymentModel),
		Labels:        ext.Map(source.Labels(), ext.UuidToString),
		FamilyMembers: ext.Map(source.FamilyMembers(), ext.UuidToString),
		Payer:         option.Map[uuid.UUID, string](source.Payer(), ext.UuidToString).Value(),
		CreatedAt:     source.CreatedAt(),
		UpdatedAt:     source.UpdatedAt(),
	}
}

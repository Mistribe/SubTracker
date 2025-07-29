package endpoints

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"
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

type subscriptionModel struct {
	Id                  string    `json:"id" binding:"required"`
	Labels              []string  `json:"labels" binding:"required"`
	FamilyMembers       []string  `json:"family_members" binding:"required"`
	PayerId             *string   `json:"payer_id_id,omitempty"`
	FamilyId            *string   `json:"family_id,omitempty"`
	PayedByJointAccount bool      `json:"payed_by_joint_account" binding:"required"`
	CreatedAt           time.Time `json:"created_at" binding:"required" format:"date-time"`
	UpdatedAt           time.Time `json:"updated_at" binding:"required" format:"date-time"`
	Etag                string    `json:"etag" binding:"required"`
}

func newSubscriptionModel(source subscription.Subscription) subscriptionModel {
	return subscriptionModel{
		Id:                  source.Id().String(),
		Name:                source.Name(),
		Payments:            slicesx.Map(source.Payments().Values(), newPaymentModel),
		Labels:              slicesx.Map(source.Labels().Values(), ext.UuidToString),
		FamilyMembers:       slicesx.Map(source.FamilyMembers().Values(), ext.UuidToString),
		PayerId:             option.Map[uuid.UUID, string](source.Payer(), ext.UuidToString).Value(),
		FamilyId:            option.Map[uuid.UUID, string](source.FamilyId(), ext.UuidToString).Value(),
		PayedByJointAccount: source.PayedByJointAccount(),
		CreatedAt:           source.CreatedAt(),
		UpdatedAt:           source.UpdatedAt(),
		Etag:                source.ETag(),
	}
}

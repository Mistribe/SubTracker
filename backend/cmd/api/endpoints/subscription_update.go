package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/pkg/ext"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]
}

type updateSubscriptionModel struct {
	FriendlyName      *string                         `json:"friendly_name,omitempty"`
	FreeTrialDays     *uint                           `json:"free_trial_days,omitempty"`
	ServiceProviderId string                          `json:"service_provider_id" binding:"required"`
	PlanId            string                          `json:"plan_id" binding:"required"`
	PriceId           string                          `json:"price_id" binding:"required"`
	ServiceUsers      []string                        `json:"service_users,omitempty"`
	StartDate         time.Time                       `json:"start_date" binding:"required" format:"date-time"`
	EndDate           *time.Time                      `json:"end_date,omitempty" format:"date-time"`
	Recurrency        string                          `json:"recurrency" binding:"required"`
	CustomRecurrency  *uint                           `json:"custom_recurrency,omitempty"`
	Payer             *editableSubscriptionPayerModel `json:"payer,omitempty"`
	Owner             editableOwnerModel              `json:"owner" binding:"required"`
	UpdatedAt         *time.Time                      `json:"updated_at,omitempty" format:"date-time"`
}

func (m updateSubscriptionModel) Subscription(userId string, subId uuid.UUID) (subscription.Subscription, error) {
	serviceProviderId, err := uuid.Parse(m.ServiceProviderId)
	if err != nil {
		return nil, err
	}
	planId, err := uuid.Parse(m.PlanId)
	if err != nil {
		return nil, err
	}
	priceId, err := uuid.Parse(m.PriceId)
	if err != nil {
		return nil, err
	}
	ownerType, err := auth.ParseOwnerType(m.Owner.Type)
	if err != nil {
		return nil, err
	}
	var familyId *uuid.UUID
	if m.Owner.FamilyId != nil {
		fid, err := uuid.Parse(*m.Owner.FamilyId)
		if err != nil {
			return nil, err
		}
		familyId = &fid
	}
	owner := auth.NewOwner(ownerType, familyId, &userId)
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	var payer subscription.Payer
	if m.Payer != nil {
		if familyId == nil {
			return nil, errors.New("missing family_id for adding a payer")
		}
		payerType, err := subscription.ParsePayerType(m.Payer.Type)
		if err != nil {
			return nil, err
		}
		var memberId *uuid.UUID
		if m.Payer.MemberId != nil {
			mbrId, err := uuid.Parse(*m.Payer.MemberId)
			if err != nil {
				return nil, err
			}
			memberId = &mbrId
		}
		payer = subscription.NewPayer(payerType, *familyId, memberId)
	}
	recurrency, err := subscription.ParseRecurrencyType(m.Recurrency)
	if err != nil {
		return nil, err
	}
	serviceUsers, err := slicesx.MapErr(m.ServiceUsers, func(in string) (uuid.UUID, error) {
		return uuid.Parse(in)
	})
	if err != nil {
		return nil, err
	}
	return subscription.NewSubscription(
		subId,
		m.FriendlyName,
		m.FreeTrialDays,
		serviceProviderId,
		planId,
		priceId,
		owner,
		payer,
		serviceUsers,
		m.StartDate,
		m.EndDate,
		recurrency,
		m.CustomRecurrency,
		updatedAt,
		updatedAt,
	), nil
}

func (m updateSubscriptionModel) Command(userId string, id uuid.UUID) (command.UpdateSubscriptionCommand, error) {
	sub, err := m.Subscription(userId, id)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	return command.UpdateSubscriptionCommand{
		Subscription: sub,
	}, nil
}

// Handle godoc
//
//	@Summary		Update an existing subscription
//	@Description	Update an existing subscription
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscriptionId	path		string					true	"Subscription ID"
//	@Param			subscription	body		updateSubscriptionModel	true	"Subscription data"
//	@Success		200				{object}	subscriptionModel
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions/{subscriptionId} [put]
func (s SubscriptionUpdateEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("subscriptionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	var model updateSubscriptionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	r := s.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return newSubscriptionModel(sub)
		}))
}

func (s SubscriptionUpdateEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (s SubscriptionUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (s SubscriptionUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionUpdateEndpoint(handler core.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]) *SubscriptionUpdateEndpoint {
	return &SubscriptionUpdateEndpoint{
		handler: handler,
	}
}

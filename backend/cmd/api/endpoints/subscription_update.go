package endpoints

import (
	"errors"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/ext"
	"net/http"
	"time"

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
	editableSubscriptionModel
	PayerType     *string    `json:"payer_type,omitempty"`
	PayerMemberId *string    `json:"payer_memberId,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" format:"date-time"`
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
	ownerType, err := user.ParseOwnerType(m.OwnerType)
	if err != nil {
		return nil, err
	}
	var familyId *uuid.UUID
	if m.FamilyId != nil {
		fid, err := uuid.Parse(*m.FamilyId)
		if err != nil {
			return nil, err
		}
		familyId = &fid
	}
	owner := user.NewOwner(ownerType, familyId, &userId)
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	var payer subscription.Payer
	if m.PayerType != nil {
		if familyId == nil {
			return nil, errors.New("missing family_id for adding a payer")
		}
		payerType, err := subscription.ParsePayerType(*m.PayerType)
		if err != nil {
			return nil, err
		}
		var memberId *uuid.UUID
		if m.PayerMemberId != nil {
			mbrId, err := uuid.Parse(*m.PayerMemberId)
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
//	@Param			id				path		string					true	"Subscription ID"
//	@Param			subscription	body		updateSubscriptionModel	true	"Subscription data"
//	@Success		200				{object}	subscriptionModel
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions/{id} [put]
func (s SubscriptionUpdateEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
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

	userId, ok := user.FromContext(c)
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
		"/:id",
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

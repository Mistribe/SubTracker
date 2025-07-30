package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionCreateEndpoint struct {
	handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionCreateEndpoint(handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]) *SubscriptionCreateEndpoint {
	return &SubscriptionCreateEndpoint{handler: handler}
}

type createSubscriptionModel struct {
	Id                *string                         `json:"id,omitempty"`
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
	CreatedAt         *time.Time                      `json:"created_at,omitempty"`
}

func (m createSubscriptionModel) Subscription(userId string) (subscription.Subscription, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
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
	ownerType, err := user.ParseOwnerType(m.Owner.Type)
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
	owner := user.NewOwner(ownerType, familyId, &userId)
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
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
		id,
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
		createdAt,
		createdAt,
	), nil
}

func (m createSubscriptionModel) Command(userId string) (command.CreateSubscriptionCommand, error) {
	sub, err := m.Subscription(userId)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	return command.CreateSubscriptionCommand{
		Subscription: sub,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new subscription
//	@Description	Create a new subscription
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		createSubscriptionModel	true	"Subscription data"
//	@Success		201				{object}	subscriptionModel
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions [post]
func (s SubscriptionCreateEndpoint) Handle(c *gin.Context) {
	var model createSubscriptionModel
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

	cmd, err := model.Command(userId)
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
		withStatus[subscription.Subscription](http.StatusCreated),
		withMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return newSubscriptionModel(sub)
		}))

}

func (s SubscriptionCreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (s SubscriptionCreateEndpoint) Method() string {
	return http.MethodPost
}

func (s SubscriptionCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

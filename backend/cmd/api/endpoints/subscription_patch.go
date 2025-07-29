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

type SubscriptionPatchEndpoint struct {
	handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionPatchEndpoint(handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]) *SubscriptionPatchEndpoint {
	return &SubscriptionPatchEndpoint{handler: handler}
}

type patchSubscriptionModel struct {
	editableSubscriptionModel
	Id            *string    `json:"id,omitempty"`
	PayerType     *string    `json:"payer_type,omitempty"`
	PayerMemberId *string    `json:"payer_memberId,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchSubscriptionModel) Subscription(userId string) (subscription.Subscription, error) {
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
		updatedAt,
		updatedAt,
	), nil
}

func (m patchSubscriptionModel) Command(userId string) (command.PatchSubscriptionCommand, error) {
	sub, err := m.Subscription(userId)
	if err != nil {
		return command.PatchSubscriptionCommand{}, nil
	}
	return command.PatchSubscriptionCommand{
		Subscription: sub,
	}, nil
}

// Handle godoc
//
//	@Summary		Update subscription
//	@Description	Update an existing subscription with new details
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		patchSubscriptionModel	true	"Subscription details to update"
//	@Success		200				{object}	subscriptionModel		"Updated subscription"
//	@Failure		400				{object}	httpError				"Invalid request"
//	@Router			/subscriptions [patch]
func (e SubscriptionPatchEndpoint) Handle(c *gin.Context) {
	var model patchSubscriptionModel
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
	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[subscription.Subscription](http.StatusOK),
		withMapping[subscription.Subscription](func(f subscription.Subscription) any {
			return newSubscriptionModel(f)
		}))
}

func (e SubscriptionPatchEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e SubscriptionPatchEndpoint) Method() string {
	return http.MethodPatch
}

func (e SubscriptionPatchEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"

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

type PatchSubscriptionModel struct {
	Id               *string                         `json:"id,omitempty"`
	FriendlyName     *string                         `json:"friendly_name,omitempty"`
	FreeTrial        *SubscriptionFreeTrialModel     `json:"free_trial,omitempty"`
	ProviderId       string                          `json:"provider_id" binding:"required"`
	PlanId           *string                         `json:"plan_id,omitempty"`
	PriceId          *string                         `json:"price_id,omitempty"`
	CustomPrice      *SubscriptionCustomPriceModel   `json:"custom_price,omitempty"`
	ServiceUsers     []string                        `json:"service_users,omitempty"`
	StartDate        time.Time                       `json:"start_date" binding:"required" format:"date-time"`
	EndDate          *time.Time                      `json:"end_date,omitempty" format:"date-time"`
	Recurrency       string                          `json:"recurrency" binding:"required"`
	CustomRecurrency *int32                          `json:"custom_recurrency,omitempty"`
	Payer            *EditableSubscriptionPayerModel `json:"payer,omitempty"`
	Owner            EditableOwnerModel              `json:"owner" binding:"required"`
	UpdatedAt        *time.Time                      `json:"updated_at,omitempty" format:"date-time"`
}

func (m PatchSubscriptionModel) Subscription(userId string) (subscription.Subscription, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	serviceProviderId, err := uuid.Parse(m.ProviderId)
	if err != nil {
		return nil, err
	}
	planId, err := parseUuidOrNil(m.PlanId)
	if err != nil {
		return nil, err
	}
	priceId, err := parseUuidOrNil(m.PriceId)
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
	serviceUsers, err := slicesx.SelectErr(m.ServiceUsers, func(in string) (uuid.UUID, error) {
		return uuid.Parse(in)
	})
	if err != nil {
		return nil, err
	}
	price, err := newSubscriptionCustomPrice(m.CustomPrice)
	if err != nil {
		return nil, err
	}

	return subscription.NewSubscription(
		id,
		m.FriendlyName,
		newSubscriptionFreeTrial(m.FreeTrial),
		serviceProviderId,
		planId,
		priceId,
		price,
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

func (m PatchSubscriptionModel) Command(userId string) (command.PatchSubscriptionCommand, error) {
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
//	@Summary		Patch subscription
//	@Description	Update or create a subscription with complete details. If subscription doesn't exist, it will be created.
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		PatchSubscriptionModel	true	"Complete subscription data"
//	@Success		200				{object}	SubscriptionModel		"Successfully updated subscription"
//	@Failure		400				{object}	httpError				"Bad Request - Invalid input data"
//	@Failure		401				{object}	httpError				"Unauthorized - Invalid user authentication"
//	@Failure		404				{object}	httpError				"Subscription not found"
//	@Failure		500				{object}	httpError				"Internal Server Error"
//	@Router			/subscriptions [patch]
func (e SubscriptionPatchEndpoint) Handle(c *gin.Context) {
	var model PatchSubscriptionModel
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

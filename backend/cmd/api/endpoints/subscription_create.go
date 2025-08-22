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

type SubscriptionCreateEndpoint struct {
	handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionCreateEndpoint(handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]) *SubscriptionCreateEndpoint {
	return &SubscriptionCreateEndpoint{handler: handler}
}

type CreateSubscriptionModel struct {
	Id               *string                         `json:"id,omitempty"`
	FriendlyName     *string                         `json:"friendly_name,omitempty"`
	FreeTrial        *SubscriptionFreeTrialModel     `json:"free_trial,omitempty"`
	ProviderId       string                          `json:"provider_id" binding:"required"`
	PlanId           *string                         `json:"plan_id,omitempty"`
	PriceId          *string                         `json:"price_id,omitempty"`
	CustomPrice      *AmountModel                    `json:"custom_price,omitempty"`
	ServiceUsers     []string                        `json:"service_users,omitempty"`
	Labels           []string                        `json:"labels,omitempty"`
	StartDate        time.Time                       `json:"start_date" binding:"required" format:"date-time"`
	EndDate          *time.Time                      `json:"end_date,omitempty" format:"date-time"`
	Recurrency       string                          `json:"recurrency" binding:"required"`
	CustomRecurrency *int32                          `json:"custom_recurrency,omitempty"`
	Payer            *EditableSubscriptionPayerModel `json:"payer,omitempty"`
	Owner            EditableOwnerModel              `json:"owner" binding:"required"`
	CreatedAt        *time.Time                      `json:"created_at,omitempty"`
}

func (m CreateSubscriptionModel) Subscription(userId string) (subscription.Subscription, error) {
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
	serviceUsers, err := slicesx.SelectErr(m.ServiceUsers, uuid.Parse)
	if err != nil {
		return nil, err
	}
	labels, err := slicesx.SelectErr(m.Labels, func(in string) (subscription.LabelRef, error) {
		labelId, err := uuid.Parse(in)
		if err != nil {
			return subscription.LabelRef{}, err
		}
		return subscription.LabelRef{
			LabelId: labelId,
			Source:  subscription.LabelSourceSubscription,
		}, nil
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
		labels,
		m.StartDate,
		m.EndDate,
		recurrency,
		m.CustomRecurrency,
		createdAt,
		createdAt,
	), nil
}

func (m CreateSubscriptionModel) Command(userId string) (command.CreateSubscriptionCommand, error) {
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
//	@Description	Create a new subscription with provider, plan, pricing, and payment information
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		CreateSubscriptionModel	true	"Subscription creation data"
//	@Success		201				{object}	SubscriptionModel		"Successfully created subscription"
//	@Failure		400				{object}	HttpErrorResponse		"Bad Request - Invalid input data"
//	@Failure		401				{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		500				{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/subscriptions [post]
func (s SubscriptionCreateEndpoint) Handle(c *gin.Context) {
	var model CreateSubscriptionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
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

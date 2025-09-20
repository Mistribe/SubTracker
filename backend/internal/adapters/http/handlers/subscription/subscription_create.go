package subscription

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/ports"
	auth2 "github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
	"github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/subscription"
)

type CreateEndpoint struct {
	handler ports.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]) *CreateEndpoint {
	return &CreateEndpoint{handler: handler}
}

func createSubscriptionRequestToSubscription(r dto.CreateSubscriptionRequest, userId string) (
	subscription.Subscription,
	error) {
	id, err := x.ParseOrNewUUID(r.Id)
	if err != nil {
		return nil, err
	}
	serviceProviderId, err := uuid.Parse(r.ProviderId)
	if err != nil {
		return nil, err
	}
	planId, err := x.ParseOrNilUUID(r.PlanId)
	if err != nil {
		return nil, err
	}
	priceId, err := x.ParseOrNilUUID(r.PriceId)
	if err != nil {
		return nil, err
	}
	ownerType, err := auth.ParseOwnerType(r.Owner.Type)
	if err != nil {
		return nil, err
	}
	var familyId *uuid.UUID
	if r.Owner.FamilyId != nil {
		fid, err := uuid.Parse(*r.Owner.FamilyId)
		if err != nil {
			return nil, err
		}
		familyId = &fid
	}
	owner := auth.NewOwner(ownerType, familyId, &userId)
	createdAt := x.ValueOrDefault(r.CreatedAt, time.Now())
	var payer subscription.Payer
	if r.Payer != nil {
		if familyId == nil {
			return nil, errors.New("missing family_id for adding a payer")
		}
		payerType, err := subscription.ParsePayerType(r.Payer.Type)
		if err != nil {
			return nil, err
		}
		var memberId *uuid.UUID
		if r.Payer.MemberId != nil {
			mbrId, err := uuid.Parse(*r.Payer.MemberId)
			if err != nil {
				return nil, err
			}
			memberId = &mbrId
		}
		payer = subscription.NewPayer(payerType, *familyId, memberId)
	}
	recurrency, err := subscription.ParseRecurrencyType(r.Recurrency)
	if err != nil {
		return nil, err
	}
	serviceUsers, err := collection.SelectErr(r.ServiceUsers, uuid.Parse)
	if err != nil {
		return nil, err
	}
	labels, err := collection.SelectErr(r.Labels, func(in string) (subscription.LabelRef, error) {
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
	price, err := dto.NewSubscriptionCustomPrice(r.CustomPrice)
	if err != nil {
		return nil, err
	}
	return subscription.NewSubscription(
		id,
		r.FriendlyName,
		dto.NewSubscriptionFreeTrial(r.FreeTrial),
		serviceProviderId,
		planId,
		priceId,
		price,
		owner,
		payer,
		serviceUsers,
		labels,
		r.StartDate,
		r.EndDate,
		recurrency,
		r.CustomRecurrency,
		createdAt,
		createdAt,
	), nil
}

// Handle godoc
//
//	@Summary		Create a new subscription
//	@Description	Create a new subscription with provider, plan, pricing, and payment information
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		dto.CreateSubscriptionRequest	true	"Subscription creation data"
//	@Success		201				{object}	dto.SubscriptionModel			"Successfully created subscription"
//	@Failure		400				{object}	HttpErrorResponse				"Bad Request - Invalid input data"
//	@Failure		401				{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		500				{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/subscriptions [post]
func (s CreateEndpoint) Handle(c *gin.Context) {
	var model dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth2.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ginx.HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	sub, err := createSubscriptionRequestToSubscription(model, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	cmd := command.CreateSubscriptionCommand{
		Subscription: sub,
	}
	r := s.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[subscription.Subscription](http.StatusCreated),
		WithMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return dto.NewSubscriptionModel(sub)
		}))

}

func (s CreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (s CreateEndpoint) Method() string {
	return http.MethodPost
}

func (s CreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

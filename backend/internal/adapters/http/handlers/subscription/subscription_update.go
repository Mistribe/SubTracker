package subscription

import (
	"errors"
	"net/http"
	"time"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/ports"
	auth2 "github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
	"github.com/mistribe/subtracker/pkg/ginx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/subscription"
)

type UpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]
}

func updateSubscriptionRequestToSubscription(
	r dto.UpdateSubscriptionRequest,
	userId string,
	subId uuid.UUID) (subscription.Subscription, error) {
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
	ownerType, err := types.ParseOwnerType(r.Owner.Type)
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
	owner := types.NewOwner(ownerType, familyId, &userId)
	updatedAt := x.ValueOrDefault(r.UpdatedAt, time.Now())
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
	serviceUsers, err := collection.SelectErr(r.ServiceUsers, func(in string) (uuid.UUID, error) {
		return uuid.Parse(in)
	})
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
		subId,
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
		updatedAt,
		updatedAt,
	), nil
}

// Handle godoc
//
//	@Summary		Update subscription by LabelID
//	@Description	Update an existing subscription's details including provider, plan, pricing, and payment information
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscriptionId	path		string							true	"Subscription LabelID (UUID format)"
//	@Param			subscription	body		dto.UpdateSubscriptionRequest	true	"Updated subscription data"
//	@Success		200				{object}	dto.SubscriptionModel			"Successfully updated subscription"
//	@Failure		400				{object}	HttpErrorResponse				"Bad Request - Invalid input data or subscription LabelID"
//	@Failure		401				{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		404				{object}	HttpErrorResponse				"Subscription not found"
//	@Failure		500				{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/subscriptions/{subscriptionId} [put]
func (s UpdateEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("subscriptionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var model dto.UpdateSubscriptionRequest
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

	sub, err := updateSubscriptionRequestToSubscription(model, userId, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	cmd := command.UpdateSubscriptionCommand{
		Subscription: sub,
	}
	r := s.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return dto.NewSubscriptionModel(sub)
		}))
}

func (s UpdateEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (s UpdateEndpoint) Method() string {
	return http.MethodPut
}

func (s UpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]) *UpdateEndpoint {
	return &UpdateEndpoint{
		handler: handler,
	}
}

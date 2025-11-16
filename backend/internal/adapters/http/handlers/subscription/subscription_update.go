package subscription

import (
	"errors"
	"net/http"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/x/herd"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/subscription"
)

type UpdateEndpoint struct {
	handler        ports.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]
	authentication ports.Authentication
}

func updateSubscriptionRequestToCommand(
	r dto.UpdateSubscriptionRequest,
	userId types.UserID,
	subscriptionID types.SubscriptionID) (command.UpdateSubscriptionCommand, error) {
	providerID, err := types.ParseProviderID(r.ProviderId)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	owner, err := types.TryParseOwnerType(r.Owner)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	var payerType *subscription.PayerType
	var payerMemberId *types.FamilyMemberID
	if r.Payer != nil {
		if owner != types.FamilyOwnerType {
			return command.UpdateSubscriptionCommand{}, errors.New("payer must be set for a family")
		}
		parsedPayerType, err := subscription.ParsePayerType(r.Payer.Type)
		if err != nil {
			return command.UpdateSubscriptionCommand{}, err
		}
		payerType = &parsedPayerType

		memberId, err := types.ParseFamilyMemberIDOrNil(r.Payer.MemberId)
		if err != nil {
			return command.UpdateSubscriptionCommand{}, err
		}
		payerMemberId = memberId
	}
	var freeTrial subscription.FreeTrial
	if r.FreeTrial != nil {
		freeTrial = subscription.NewFreeTrial(r.FreeTrial.StartDate, r.FreeTrial.EndDate)
	}
	recurrency, err := subscription.ParseRecurrencyType(r.Recurrency)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	familyUsers, err := herd.SelectErr(r.ServiceUsers, types.ParseFamilyMemberID)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	labels, err := herd.SelectErr(r.Labels, func(in string) (types.LabelID, error) {
		labelId, err := types.ParseLabelID(in)
		if err != nil {
			return types.LabelID{}, err
		}
		return labelId, nil
	})
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	price, err := dto.NewSubscriptionCustomPrice(r.CustomPrice)
	if err != nil {
		return command.UpdateSubscriptionCommand{}, err
	}
	return command.UpdateSubscriptionCommand{
		SubscriptionID:   subscriptionID,
		FriendlyName:     r.FriendlyName,
		FreeTrial:        freeTrial,
		ProviderID:       providerID,
		Price:            price.Amount(),
		Owner:            owner,
		PayerType:        payerType,
		PayerMemberId:    payerMemberId,
		FamilyUsers:      familyUsers,
		Labels:           labels,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		Recurrency:       recurrency,
		CustomRecurrency: r.CustomRecurrency,
		UpdatedAt:        option.New(r.UpdatedAt),
	}, nil
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
	subscriptionID, err := types.ParseSubscriptionID(c.Param("subscriptionId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var model dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	connectedAccount := s.authentication.MustGetConnectedAccount(c)
	cmd, err := updateSubscriptionRequestToCommand(model, connectedAccount.UserID(), subscriptionID)
	if err != nil {
		FromError(c, err)
		return
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

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription],
	authentication ports.Authentication) *UpdateEndpoint {
	return &UpdateEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}

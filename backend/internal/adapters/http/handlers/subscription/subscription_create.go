package subscription

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/x/herd"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

type CreateEndpoint struct {
	handler        ports.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
	authentication ports.Authentication
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription],
	authentication ports.Authentication) *CreateEndpoint {
	return &CreateEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}

func createSubscriptionRequestToCommand(r dto.CreateSubscriptionRequest) (command.CreateSubscriptionCommand, error) {
	subscriptionID, err := types.ParseSubscriptionIDOrNil(r.Id)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	providerID, err := types.ParseProviderIDOrNil(r.ProviderId)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	owner, err := types.TryParseOwnerType(r.Owner)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	var payerType *subscription.PayerType
	var payerMemberId *types.FamilyMemberID
	if r.Payer != nil {
		if owner != types.FamilyOwnerType {
			return command.CreateSubscriptionCommand{}, errors.New("payer must be set for a family")
		}
		parsedPayerType, err := subscription.ParsePayerType(r.Payer.Type)
		if err != nil {
			return command.CreateSubscriptionCommand{}, err
		}
		payerType = &parsedPayerType

		memberId, err := types.ParseFamilyMemberIDOrNil(r.Payer.MemberId)
		if err != nil {
			return command.CreateSubscriptionCommand{}, err
		}
		payerMemberId = memberId
	}
	var freeTrial subscription.FreeTrial
	if r.FreeTrial != nil {
		freeTrial = subscription.NewFreeTrial(r.FreeTrial.StartDate, r.FreeTrial.EndDate)
	}
	recurrency, err := subscription.ParseRecurrencyType(r.Recurrency)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	familyUsers, err := herd.SelectErr(r.FamilyUsers, types.ParseFamilyMemberID)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	labels, err := herd.SelectErr(r.Labels, func(in string) (types.LabelID, error) {
		labelId, err := types.ParseLabelID(in)
		if err != nil {
			return types.LabelID{}, err
		}
		return labelId, nil
	})
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	price, err := dto.NewSubscriptionCustomPrice(r.Price)
	if err != nil {
		return command.CreateSubscriptionCommand{}, err
	}
	if r.ProviderKey == nil && providerID == nil {
		return command.CreateSubscriptionCommand{}, errors.New("provider id or provider key is required")
	}
	return command.CreateSubscriptionCommand{
		SubscriptionID:   option.New(subscriptionID),
		FriendlyName:     r.FriendlyName,
		FreeTrial:        freeTrial,
		ProviderID:       providerID,
		ProviderKey:      r.ProviderKey,
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
		CreatedAt:        option.New(r.CreatedAt),
	}, nil
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
		FromError(c, err)
		return
	}

	cmd, err := createSubscriptionRequestToCommand(model)
	if err != nil {
		FromError(c, err)
		return
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

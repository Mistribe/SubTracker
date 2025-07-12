package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type SubscriptionPaymentCreateEndpoint struct {
	handler core.CommandHandler[command.CreatePaymentCommand, subscription.Subscription]
}

func NewSubscriptionPaymentCreateEndpoint(handler core.CommandHandler[command.CreatePaymentCommand, subscription.Subscription]) *SubscriptionPaymentCreateEndpoint {
	return &SubscriptionPaymentCreateEndpoint{handler: handler}
}

type createSubscriptionPaymentModel struct {
	Id        *string    `json:"id,omitempty"`
	Price     float64    `json:"price"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Months    int        `json:"months"`
	Currency  string     `json:"currency"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m createSubscriptionPaymentModel) ToPayment() result.Result[subscription.Payment] {
	var id uuid.UUID
	var err error
	var endDate option.Option[time.Time]
	var createdAt time.Time

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[subscription.Payment](err)
	}
	endDate = option.New(m.EndDate)
	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())
	return subscription.NewPayment(
		id,
		m.Price,
		m.StartDate,
		endDate,
		m.Months,
		m.Currency,
		createdAt,
		createdAt,
	)
}

func (m createSubscriptionPaymentModel) Command(subscriptionId string) result.Result[command.CreatePaymentCommand] {
	return result.Bind[subscription.Payment, command.CreatePaymentCommand](
		m.ToPayment(),
		func(payment subscription.Payment) result.Result[command.CreatePaymentCommand] {
			subId, err := uuid.Parse(subscriptionId)
			if err != nil {
				return result.Fail[command.CreatePaymentCommand](err)
			}
			return result.Success(command.CreatePaymentCommand{
				SubscriptionId: subId,
				Payment:        payment,
			})
		})
}

// Handle godoc
//
//	@Summary		Create a new subscription payment
//	@Description	Create a new subscription payment
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription_id	path		string							true	"Subscription ID"
//	@Param			payment			body		createSubscriptionPaymentModel	true	"Payment data"
//	@Success		201				{object}	subscriptionModel
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions/{subscription_id}/payments [post]
func (s SubscriptionPaymentCreateEndpoint) Handle(c *gin.Context) {
	var model createSubscriptionPaymentModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	subscriptionId := c.Param("id")
	result.Match[command.CreatePaymentCommand, result.Unit](model.Command(subscriptionId),
		func(cmd command.CreatePaymentCommand) result.Unit {
			r := s.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[subscription.Subscription](http.StatusCreated),
				withMapping[subscription.Subscription](func(sub subscription.Subscription) any {
					return newSubscriptionModel(sub)
				}))
			return result.Unit{}
		},
		func(err error) result.Unit {
			c.JSON(http.StatusBadRequest, httpError{
				Message: err.Error(),
			})
			return result.Unit{}
		})
}

func (s SubscriptionPaymentCreateEndpoint) Pattern() []string {
	return []string{
		"/:id/payments",
	}
}

func (s SubscriptionPaymentCreateEndpoint) Method() string {
	return http.MethodPost
}

func (s SubscriptionPaymentCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

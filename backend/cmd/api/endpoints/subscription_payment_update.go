package endpoints

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type SubscriptionPaymentUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdatePaymentCommand, subscription.Subscription]
}

type updatePaymentModel struct {
	Price     float64    `json:"price"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Months    int        `json:"months"`
	Currency  string     `json:"currency"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m updatePaymentModel) ToCommand(
	subscriptionId uuid.UUID,
	paymentId uuid.UUID) result.Result[command.UpdatePaymentCommand] {
	endDate := option.New(m.EndDate)
	updatedAt := option.New(m.UpdatedAt)
	paymentCurrency, err := currency.ParseISO(strings.TrimSpace(m.Currency))
	if err != nil {
		return result.Fail[command.UpdatePaymentCommand](err)
	}
	return result.Success(command.UpdatePaymentCommand{
		SubscriptionId: subscriptionId,
		PaymentId:      paymentId,
		Price:          m.Price,
		StartDate:      m.StartDate,
		EndDate:        endDate,
		Months:         m.Months,
		Currency:       paymentCurrency,
		UpdatedAt:      updatedAt,
	})
}

// Handle godoc
// @Summary      Update subscription payment
// @Description  Update payment details for a specific subscription
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        id        path     string              true  "Subscription ID"
// @Param        paymentId path     string              true  "Payment ID"
// @Param        payment   body     updatePaymentModel  true  "Payment details"
// @Success      200       {object} subscriptionModel
// @Failure      400       {object} httpError
// @Failure      404       {object} httpError
// @Router       /subscriptions/{id}/payments/{paymentId} [put]
func (e SubscriptionPaymentUpdateEndpoint) Handle(c *gin.Context) {
	subscriptionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
		return
	}

	paymentId, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
		return
	}

	var model updatePaymentModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
		return
	}

	result.Match(model.ToCommand(subscriptionId, paymentId),
		func(cmd command.UpdatePaymentCommand) result.Unit {
			r := e.handler.Handle(c, cmd)
			handleResponse(c, r,
				withMapping[subscription.Subscription](func(sub subscription.Subscription) any {
					return newSubscriptionModel(sub)
				}))
			return result.Unit{}
		},
		func(err error) result.Unit {
			c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
			return result.Unit{}
		})
}

func (e SubscriptionPaymentUpdateEndpoint) Pattern() []string {
	return []string{
		"/:id/payments/:paymentId",
	}
}

func (e SubscriptionPaymentUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e SubscriptionPaymentUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionPaymentUpdateEndpoint(handler core.CommandHandler[command.UpdatePaymentCommand, subscription.Subscription]) *SubscriptionPaymentUpdateEndpoint {
	return &SubscriptionPaymentUpdateEndpoint{
		handler: handler,
	}
}

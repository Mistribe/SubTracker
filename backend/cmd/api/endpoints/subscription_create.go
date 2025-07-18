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

type SubscriptionCreateEndpoint struct {
	handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionCreateEndpoint(handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]) *SubscriptionCreateEndpoint {
	return &SubscriptionCreateEndpoint{handler: handler}
}

type createPaymentModel struct {
	Id        *string    `json:"id,omitempty"`
	Price     float64    `json:"price"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Months    int        `json:"months"`
	Currency  string     `json:"currency"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m createPaymentModel) ToPayment() result.Result[subscription.Payment] {
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

type createSubscriptionModel struct {
	Id                  *string              `json:"id,omitempty"`
	FamilyId            *string              `json:"family_id,omitempty"`
	Name                string               `json:"name"`
	Payments            []createPaymentModel `json:"payments"`
	Labels              []string             `json:"labels"`
	FamilyMembers       []string             `json:"family_members"`
	PayerId             *string              `json:"payer_id,omitempty"`
	PayedByJointAccount bool                 `json:"payed_by_joint_account,omitempty"`
	CreatedAt           *time.Time           `json:"created_at,omitempty"`
}

func (m createSubscriptionModel) ToSubscription() result.Result[subscription.Subscription] {
	var id uuid.UUID
	var labels []uuid.UUID
	var familyMembers []uuid.UUID
	var payerId option.Option[uuid.UUID]
	var createdAt time.Time
	var err error
	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	paymentRes := ext.Map(m.Payments, func(in createPaymentModel) result.Result[subscription.Payment] {
		return in.ToPayment()
	})
	labels, err = ext.MapErr(m.Labels, uuid.Parse)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	familyMembers, err = ext.MapErr(m.FamilyMembers, uuid.Parse)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	payerId, err = option.ParseNew(m.PayerId, uuid.Parse)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())
	familyId, err := option.ParseNew(m.FamilyId, uuid.Parse)

	return result.BindReduce(paymentRes, func(value []subscription.Payment) result.Result[subscription.Subscription] {
		return subscription.NewSubscription(
			id,
			familyId,
			m.Name,
			value,
			labels,
			familyMembers,
			payerId,
			m.PayedByJointAccount,
			createdAt,
			createdAt,
		)
	})

}

func (m createSubscriptionModel) Command() result.Result[command.CreateSubscriptionCommand] {
	return result.Bind[subscription.Subscription, command.CreateSubscriptionCommand](
		m.ToSubscription(),
		func(sub subscription.Subscription) result.Result[command.CreateSubscriptionCommand] {
			return result.Success(command.CreateSubscriptionCommand{
				Subscription: sub,
			})
		})
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

	result.Match[command.CreateSubscriptionCommand, result.Unit](model.Command(),
		func(cmd command.CreateSubscriptionCommand) result.Unit {
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

package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
	"net/http"
	"time"
)

type SubscriptionCreateEndpoint struct {
	handler core.CommandHandler[command.CreateSubscriptionCommand, subscription.Subscription]
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
	Id            *string              `json:"id,omitempty"`
	Name          string               `json:"name"`
	Payments      []createPaymentModel `json:"payments"`
	Labels        []string             `json:"labels"`
	FamilyMembers []string             `json:"family_members"`
	Payer         *string              `json:"payer,omitempty"`
	CreatedAt     *time.Time           `json:"created_at,omitempty"`
}

func (m createSubscriptionModel) ToSubscription() result.Result[subscription.Subscription] {
	var id uuid.UUID
	var labels []uuid.UUID
	var familyMembers []uuid.UUID
	var payer option.Option[uuid.UUID]
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
	payer, err = option.ParseNew(m.Payer, uuid.Parse)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())

	return result.BindReduce(paymentRes, func(value []subscription.Payment) result.Result[subscription.Subscription] {
		return subscription.NewSubscription(
			id,
			m.Name,
			value,
			labels,
			familyMembers,
			payer,
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
		"/:id",
	}
}

func (s SubscriptionCreateEndpoint) Method() string {
	return http.MethodGet
}

func (s SubscriptionCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionCreateEndpoint() *SubscriptionCreateEndpoint {
	return &SubscriptionCreateEndpoint{}
}

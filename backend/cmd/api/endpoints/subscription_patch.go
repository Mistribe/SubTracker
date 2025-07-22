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
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type SubscriptionPatchEndpoint struct {
	handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionPatchEndpoint(handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]) *SubscriptionPatchEndpoint {
	return &SubscriptionPatchEndpoint{handler: handler}
}

type patchPaymentModel struct {
	Id        *string    `json:"id,omitempty"`
	Price     float64    `json:"price"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Months    int        `json:"months"`
	Currency  string     `json:"currency"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m patchPaymentModel) ToPayment(subscriptionId uuid.UUID) (subscription.Payment, error) {
	var id uuid.UUID
	var err error
	var endDate option.Option[time.Time]

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return subscription.Payment{}, err
	}
	endDate = option.New(m.EndDate)
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	paymentCurrency, err := currency.ParseISO(strings.TrimSpace(m.Currency))
	if err != nil {
		return subscription.Payment{}, err
	}
	return subscription.NewPaymentWithoutValidation(
		id,
		m.Price,
		m.StartDate,
		endDate,
		m.Months,
		paymentCurrency,
		subscriptionId,
		updatedAt,
		updatedAt,
		false,
	), nil
}

type patchSubscriptionModel struct {
	Id                  *string             `json:"id,omitempty"`
	FamilyId            *string             `json:"family_id,omitempty"`
	Name                string              `json:"name"`
	Payments            []patchPaymentModel `json:"payments"`
	Labels              []string            `json:"labels"`
	FamilyMembers       []string            `json:"family_members"`
	PayerId             *string             `json:"payer_id,omitempty"`
	PayedByJointAccount bool                `json:"payed_by_joint_account,omitempty"`
	UpdatedAt           *time.Time          `json:"updated_at,omitempty"`
}

func (m patchSubscriptionModel) Command() result.Result[command.PatchSubscriptionCommand] {
	var id uuid.UUID
	var labels []uuid.UUID
	var familyMembers []uuid.UUID
	var payerId option.Option[uuid.UUID]
	var updatedAt time.Time
	var err error
	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	labels, err = ext.MapErr(m.Labels, uuid.Parse)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	familyMembers, err = ext.MapErr(m.FamilyMembers, uuid.Parse)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	payerId, err = option.ParseNew(m.PayerId, uuid.Parse)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	updatedAt = ext.ValueOrDefault(m.UpdatedAt, time.Now())
	familyId, err := option.ParseNew(m.FamilyId, uuid.Parse)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	payments, err := ext.MapErr(m.Payments, func(value patchPaymentModel) (subscription.Payment, error) {
		return value.ToPayment(id)
	})
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}

	return result.Success(command.PatchSubscriptionCommand{
		Subscription: subscription.NewSubscriptionWithoutValidation(
			id,
			familyId.Value(),
			m.Name,
			payments,
			labels,
			familyMembers,
			payerId.Value(),
			m.PayedByJointAccount,
			updatedAt,
			updatedAt,
			false,
		),
	})
}

func (e SubscriptionPatchEndpoint) Handle(c *gin.Context) {
	var model patchSubscriptionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := model.Command()
	result.Match[command.PatchSubscriptionCommand, result.Unit](cmd,
		func(cmd command.PatchSubscriptionCommand) result.Unit {
			r := e.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[subscription.Subscription](http.StatusOK),
				withMapping[subscription.Subscription](func(f subscription.Subscription) any {
					return newSubscriptionModel(f)
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

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
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionPatchEndpoint struct {
	handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]
}

func NewSubscriptionPatchEndpoint(handler core.CommandHandler[command.PatchSubscriptionCommand, subscription.Subscription]) *SubscriptionPatchEndpoint {
	return &SubscriptionPatchEndpoint{handler: handler}
}

// @Description	Payment update model
type patchSubscriptionPaymentModel struct {
	Id        *string    `json:"id,omitempty"`
	Price     float64    `json:"price" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Months    int        `json:"months" binding:"required"`
	Currency  string     `json:"currency" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchSubscriptionPaymentModel) ToPayment(subscriptionId uuid.UUID) (subscription.Payment, error) {
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

// @Description	Subscription update model
type patchSubscriptionModel struct {
	Id                  *string                         `json:"id,omitempty"`
	FamilyId            *string                         `json:"family_id,omitempty"`
	Name                string                          `json:"name" binding:"required"`
	Payments            []patchSubscriptionPaymentModel `json:"payments" binding:"required"`
	Labels              []string                        `json:"labels" binding:"required"`
	FamilyMembers       []string                        `json:"family_members" binding:"required"`
	PayerId             *string                         `json:"payer_id,omitempty"`
	PayedByJointAccount bool                            `json:"payed_by_joint_account,omitempty"`
	UpdatedAt           *time.Time                      `json:"updated_at,omitempty" format:"date-time"`
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
	labels, err = slicesx.MapErr(m.Labels, uuid.Parse)
	if err != nil {
		return result.Fail[command.PatchSubscriptionCommand](err)
	}
	familyMembers, err = slicesx.MapErr(m.FamilyMembers, uuid.Parse)
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
	payments, err := slicesx.MapErr(m.Payments,
		func(value patchSubscriptionPaymentModel) (subscription.Payment, error) {
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

// Handle godoc
//
//	@Summary		Update subscription
//	@Description	Update an existing subscription with new details
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		patchSubscriptionModel	true	"Subscription details to update"
//	@Success		200				{object}	subscriptionModel		"Updated subscription"
//	@Failure		400				{object}	httpError				"Invalid request"
//	@Router			/subscriptions [patch]
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

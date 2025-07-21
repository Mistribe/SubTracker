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

type SubscriptionUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]
}

type updateSubscriptionModel struct {
	Name                string     `json:"name"`
	Labels              []string   `json:"labels"`
	FamilyMembers       []string   `json:"family_members"`
	FamilyId            *string    `json:"family_id,omitempty"`
	PayerId             *string    `json:"payer_id,omitempty"`
	PayedByJointAccount bool       `json:"payed_by_joint_account,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

func (m updateSubscriptionModel) Command(id uuid.UUID) result.Result[command.UpdateSubscriptionCommand] {
	var labels []uuid.UUID
	var familyMembers []uuid.UUID
	var payer option.Option[uuid.UUID]
	var err error

	labels, err = ext.MapErr(m.Labels, uuid.Parse)
	if err != nil {
		return result.Fail[command.UpdateSubscriptionCommand](err)
	}

	familyMembers, err = ext.MapErr(m.FamilyMembers, uuid.Parse)
	if err != nil {
		return result.Fail[command.UpdateSubscriptionCommand](err)
	}

	payer, err = option.ParseNew(m.PayerId, uuid.Parse)
	if err != nil {
		return result.Fail[command.UpdateSubscriptionCommand](err)
	}

	familyId, err := option.ParseNew(m.FamilyId, uuid.Parse)
	if err != nil {
		return result.Fail[command.UpdateSubscriptionCommand](err)
	}

	return result.Success(command.UpdateSubscriptionCommand{
		Id:                  id,
		Name:                m.Name,
		Labels:              labels,
		FamilyMembers:       familyMembers,
		PayerId:             payer,
		PayedByJointAccount: m.PayedByJointAccount,
		UpdatedAt:           option.New(m.UpdatedAt),
		FamilyId:            familyId,
	})
}

// Handle godoc
//
//	@Summary		Update an existing subscription
//	@Description	Update an existing subscription
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string					true	"Subscription ID"
//	@Param			subscription	body		updateSubscriptionModel	true	"Subscription data"
//	@Success		200				{object}	subscriptionModel
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions/{id} [put]
func (s SubscriptionUpdateEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	var model updateSubscriptionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	result.Match[command.UpdateSubscriptionCommand, result.Unit](model.Command(id),
		func(cmd command.UpdateSubscriptionCommand) result.Unit {
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

func (s SubscriptionUpdateEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (s SubscriptionUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (s SubscriptionUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionUpdateEndpoint(handler core.CommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription]) *SubscriptionUpdateEndpoint {
	return &SubscriptionUpdateEndpoint{
		handler: handler,
	}
}

package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/internal/application/family/query"
)

type FamilySeeInvitationEndpoint struct {
	handler core.QueryHandler[query.SeeInvitationQuery, query.SeeInvitationQueryResponse]
}

type FamilySeeInvitationRequest struct {
	// Family unique identifier
	FamilyId string `uri:"familyId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Invitation code to verify
	Code string `form:"code" binding:"required" example:"ABC123"`
	// Family member unique identifier
	FamilyMemberId string `form:"family_member_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type FamilySeeInvitationResponse struct {
	// Family details
	Family familyModel `json:"family"`
	// Role of the invited member
	InvitedInasmuchAs string `json:"invited_inasmuch_as" example:"OWNER"`
}

func NewFamilySeeInvitationEndpoint(handler core.QueryHandler[query.SeeInvitationQuery, query.SeeInvitationQueryResponse]) FamilySeeInvitationEndpoint {
	return FamilySeeInvitationEndpoint{
		handler: handler,
	}
}

// Handle godoc
//
//	@Summary		View family invitation details
//	@Description	Get information about a family invitation using invitation code
//	@Tags			Family
//	@Produce		json
//	@Param			familyId			path		string	true	"Family ID"
//	@Param			code				query		string	true	"Invitation code"
//	@Param			family_member_id	query		string	true	"Family member ID"
//	@Success		200					{object}	FamilySeeInvitationResponse
//	@Failure		400					{object}	HttpErrorResponse
//	@Router			/families/{familyId}/invitation [get]
func (e FamilySeeInvitationEndpoint) Handle(c *gin.Context) {
	var req FamilySeeInvitationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	familyId, err := uuid.Parse(req.FamilyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	familyMemberId, err := uuid.Parse(req.FamilyMemberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, query.SeeInvitationQuery{
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
		InvitationCode: req.Code,
	})
	handleResponse(c, r, withMapping[query.SeeInvitationQueryResponse](func(res query.SeeInvitationQueryResponse) any {
		return FamilySeeInvitationResponse{
			Family:            newFamilyModel(res.UserId, res.Family),
			InvitedInasmuchAs: res.InvitedInasmuchAs.String(),
		}
	}))
}

func (e FamilySeeInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/invitation",
	}
}

func (e FamilySeeInvitationEndpoint) Method() string {
	return http.MethodGet
}

func (e FamilySeeInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

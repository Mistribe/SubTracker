package family

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

type SeeInvitationEndpoint struct {
	handler ports.QueryHandler[query.SeeInvitationQuery, query.SeeInvitationQueryResponse]
}

type SeeInvitationRequest struct {
	// Family unique identifier
	FamilyId string `uri:"familyId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Invitation code to verify
	Code string `form:"code" binding:"required" example:"ABC123"`
	// Family member unique identifier
	FamilyMemberId string `form:"family_member_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

func NewSeeInvitationEndpoint(handler ports.QueryHandler[query.SeeInvitationQuery, query.SeeInvitationQueryResponse]) *SeeInvitationEndpoint {
	return &SeeInvitationEndpoint{
		handler: handler,
	}
}

// Handle godoc
//
//	@Summary		View family invitation details
//	@Description	Get information about a family invitation using invitation code
//	@Tags			family
//	@Produce		json
//	@Param			familyId			path		string	true	"Family LabelID"
//	@Param			code				query		string	true	"Invitation code"
//	@Param			family_member_id	query		string	true	"Family member LabelID"
//	@Success		200					{object}	dto.FamilySeeInvitationResponse
//	@Failure		400					{object}	HttpErrorResponse
//	@Router			/family/{familyId}/invitation [get]
func (e SeeInvitationEndpoint) Handle(c *gin.Context) {
	var req SeeInvitationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		FromError(c, err)
		return
	}

	familyId, err := uuid.Parse(req.FamilyId)
	if err != nil {
		FromError(c, err)
		return
	}

	familyMemberId, err := uuid.Parse(req.FamilyMemberId)
	if err != nil {
		FromError(c, err)
		return
	}

	r := e.handler.Handle(c, query.SeeInvitationQuery{
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
		InvitationCode: req.Code,
	})
	FromResult(c, r,
		WithMapping[query.SeeInvitationQueryResponse](func(res query.SeeInvitationQueryResponse) any {
			return dto.FamilySeeInvitationResponse{
				Family:            dto.NewFamilyModel(res.UserId, res.Family),
				InvitedInasmuchAs: res.InvitedInasmuchAs.String(),
			}
		}))
}

func (e SeeInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/invitation",
	}
}

func (e SeeInvitationEndpoint) Method() string {
	return http.MethodGet
}

func (e SeeInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

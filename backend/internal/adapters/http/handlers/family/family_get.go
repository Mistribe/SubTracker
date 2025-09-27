package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

type GetEndpoint struct {
	handler        ports.QueryHandler[query.FindUserFamilyQuery, query.FindUserFamilyQueryResponse]
	authentication ports.Authentication
}

// Handle godoc
//
//	@Summary		Get user's family
//	@Description	Retrieve the user's family
//	@Tags			family
//	@Produce		json
//	@Success		200	{object}	dto.UserFamilyResponse	"Successfully retrieved family"
//	@Failure		400	{object}	HttpErrorResponse		"Bad Request - Invalid LabelID format"
//	@Failure		401	{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		404	{object}	HttpErrorResponse		"Family not found"
//	@Failure		500	{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/family [get]
func (e GetEndpoint) Handle(c *gin.Context) {
	connectedAccount := e.authentication.MustGetConnectedAccount(c)

	q := query.FindUserFamilyQuery{
		UserId: connectedAccount.UserID(),
	}

	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[query.FindUserFamilyQueryResponse](func(r query.FindUserFamilyQueryResponse) any {
			return dto.UserFamilyResponse{
				Family: dto.NewFamilyModel(connectedAccount.UserID(), r.Family),
				Limits: dto.NewLimits(r.Limits),
			}
		}))
}

func (e GetEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e GetEndpoint) Method() string {
	return http.MethodGet
}

func (e GetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewGetEndpoint(handler ports.QueryHandler[query.FindUserFamilyQuery, query.FindUserFamilyQueryResponse],
	authentication ports.Authentication) *GetEndpoint {
	return &GetEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}

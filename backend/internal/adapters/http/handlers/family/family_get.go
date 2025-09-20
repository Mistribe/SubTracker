package family

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

type GetEndpoint struct {
	handler ports.QueryHandler[query.FindUserFamilyQuery, query.FindUserFamilyQueryResponse]
}

// Handle godoc
//
//	@Summary		Get user's family
//	@Description	Retrieve the user's family
//	@Tags			family
//	@Produce		json
//	@Success		200	{object}	dto.UserFamilyResponse	"Successfully retrieved family"
//	@Failure		400	{object}	HttpErrorResponse		"Bad Request - Invalid ID format"
//	@Failure		401	{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		404	{object}	HttpErrorResponse		"Family not found"
//	@Failure		500	{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/family [get]
func (f GetEndpoint) Handle(c *gin.Context) {
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		FromError(c, errors.New("invalid user id"))
		return
	}

	q := query.FindUserFamilyQuery{
		UserId: userId,
	}

	r := f.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[query.FindUserFamilyQueryResponse](func(r query.FindUserFamilyQueryResponse) any {
			return dto.UserFamilyResponse{
				Family: dto.NewFamilyModel(userId, r.Family),
				Limits: dto.NewLimits(r.Limits),
			}
		}))
}

func (f GetEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (f GetEndpoint) Method() string {
	return http.MethodGet
}

func (f GetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewGetEndpoint(handler ports.QueryHandler[query.FindUserFamilyQuery, query.FindUserFamilyQueryResponse]) *GetEndpoint {
	return &GetEndpoint{
		handler: handler,
	}
}

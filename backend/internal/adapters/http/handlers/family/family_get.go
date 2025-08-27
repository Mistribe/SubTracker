package family

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

type FamilyGetEndpoint struct {
	handler ports.QueryHandler[query.FindUserFamilyQuery, family.Family]
}

// Handle godoc
//
//	@Summary		Get user's family
//	@Description	Retrieve the user's family
//	@Tags			family
//	@Produce		json
//	@Success		200	{object}	familyModel			"Successfully retrieved family"
//	@Failure		400	{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		401	{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		404	{object}	HttpErrorResponse	"Family not found"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families/me [get]
func (f FamilyGetEndpoint) Handle(c *gin.Context) {
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
		WithMapping[family.Family](func(fm family.Family) any {
			return newFamilyModel(userId, fm)
		}))
}

func (f FamilyGetEndpoint) Pattern() []string {
	return []string{
		"/me",
	}
}

func (f FamilyGetEndpoint) Method() string {
	return http.MethodGet
}

func (f FamilyGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberGetEndpoint(handler ports.QueryHandler[query.FindUserFamilyQuery, family.Family]) *FamilyGetEndpoint {
	return &FamilyGetEndpoint{
		handler: handler,
	}
}

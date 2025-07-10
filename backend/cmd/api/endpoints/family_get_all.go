package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyMemberGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, []family.Member]
}

// Handle godoc
// @Summary		Get all family members
// @Description	Get all family members
// @Tags			family
// @Produce		json
// @Success		200	{array}		familyMemberModel
// @Failure		400	{object}	httpError
// @Router			/families/members [get]
func (f FamilyMemberGetAllEndpoint) Handle(c *gin.Context) {
	q := query.FindAllQuery{}

	r := f.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[[]family.Member](func(fms []family.Member) any {
			result := make([]interface{}, len(fms))
			for i, fm := range fms {
				result[i] = newFamilyMemberModel(fm)
			}
			return result
		}))
}

func (f FamilyMemberGetAllEndpoint) Pattern() []string {
	return []string{
		"/members",
	}
}

func (f FamilyMemberGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (f FamilyMemberGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, []family.Member]) *FamilyMemberGetAllEndpoint {
	return &FamilyMemberGetAllEndpoint{
		handler: handler,
	}
}

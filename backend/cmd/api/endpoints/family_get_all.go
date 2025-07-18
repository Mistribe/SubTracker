package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, []family.Family]
}

// Handle godoc
// @Summary		Get all family members
// @Description	Get all family members
// @Tags			family
// @Produce		json
// @Success		200	{array}		familyModel
// @Failure		400	{object}	httpError
// @Router			/families [get]
func (f FamilyGetAllEndpoint) Handle(c *gin.Context) {
	q := query.FindAllQuery{}

	r := f.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[[]family.Family](func(fms []family.Family) any {
			result := make([]interface{}, len(fms))
			for i, fm := range fms {
				result[i] = newFamilyModel(fm)
			}
			return result
		}))
}

func (f FamilyGetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (f FamilyGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (f FamilyGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, []family.Family]) *FamilyGetAllEndpoint {
	return &FamilyGetAllEndpoint{
		handler: handler,
	}
}

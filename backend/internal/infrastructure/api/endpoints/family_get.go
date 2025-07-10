package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyMemberGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, family.Member]
}

func (f FamilyMemberGetEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid id format",
		})
		return
	}

	q := query.FindOneQuery{
		Id: id,
	}

	r := f.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[family.Member](func(fm family.Member) any {
			return newFamilyMemberModel(fm)
		}))
}

func (f FamilyMemberGetEndpoint) Pattern() []string {
	return []string{
		"/members/:id",
	}
}

func (f FamilyMemberGetEndpoint) Method() string {
	return http.MethodGet
}

func (f FamilyMemberGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberGetEndpoint(handler core.QueryHandler[query.FindOneQuery, family.Member]) *FamilyMemberGetEndpoint {
	return &FamilyMemberGetEndpoint{
		handler: handler,
	}
}

package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[family.Family]]
}

// Handle godoc
//
//	@Summary		Get all family members
//	@Description	Get all family members
//	@Tags			family
//	@Produce		json
//	@Param			size	query		integer	false	"Number of items per page"
//	@Param			page	query		integer	false	"Page number"
//	@Success		200		{object}	paginatedResponseModel[familyModel]
//	@Failure		400		{object}	httpError
//	@Router			/families [get]
func (f FamilyGetAllEndpoint) Handle(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	q := query.NewFindAllQuery(size, page)
	userId, ok := user.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}
	r := f.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[family.Family]](func(paginatedResult core.PaginatedResponse[family.Family]) any {
			return newPaginatedResponseModel(paginatedResult, func(value family.Family) any {
				return newFamilyModel(userId, value)
			})
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

func NewFamilyMemberGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[family.Family]]) *FamilyGetAllEndpoint {
	return &FamilyGetAllEndpoint{
		handler: handler,
	}
}

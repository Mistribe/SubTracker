package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[family.Family]]
}

// Handle godoc
//
//	@Summary		Get all families
//	@Description	Retrieve a paginated list of families for the authenticated user
//	@Tags			family
//	@Produce		json
//	@Param			limit	query		integer								false	"Number of items per page (default: 10)"
//	@Param			offset	query		integer								false	"Page number (default: 1)"
//	@Success		200		{object}	PaginatedResponseModel[familyModel]	"Paginated list of families"
//	@Failure		400		{object}	httpError							"Bad Request - Invalid query parameters"
//	@Failure		401		{object}	httpError							"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	httpError							"Internal Server Error"
//	@Router			/families [get]
func (f FamilyGetAllEndpoint) Handle(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	q := query.NewFindAllQuery(limit, offset)
	userId, ok := auth.GetUserIdFromContext(c)
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

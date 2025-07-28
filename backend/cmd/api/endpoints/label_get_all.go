package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[label.label]]
}

func NewLabelGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[label.label]]) *LabelGetAllEndpoint {
	return &LabelGetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all labels
//	@Description	Get all labels
//	@Tags			label
//	@Produce		json
//	@Param			with_default	query		boolean	false	"Include default labels"
//	@Param			size			query		integer	false	"Number of items per page"
//	@Param			page			query		integer	false	"Offset number"
//	@Success		200				{object}	paginatedResponseModel[labelModel]
//	@Failure		400				{object}	httpError
//	@Router			/labels [get]
func (e LabelGetAllEndpoint) Handle(c *gin.Context) {
	withDefault := c.DefaultQuery("with_default", "false") == "true"
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	q := query.NewFindAllQuery(withDefault, size, page)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[label.label]](func(paginatedResult core.PaginatedResponse[label.label]) any {
			return newPaginatedResponseModel(paginatedResult, newLabelModel)
		}))
}

func (e LabelGetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e LabelGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (e LabelGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

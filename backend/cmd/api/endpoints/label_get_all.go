package endpoints

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[label.Label]]
}

func NewLabelGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[label.Label]]) *LabelGetAllEndpoint {
	return &LabelGetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all labels
//	@Description	Retrieve a paginated list of labels with optional filtering by owner type
//	@Tags			label
//	@Produce		json
//	@Param			owner_type	query		[]string							false	"Owner types to filter by (system,personal,family). Can be provided multiple times."
//	@Param			size		query		integer								false	"Number of items per page (default: 10)"
//	@Param			page		query		integer								false	"Page number (default: 1)"
//	@Success		200			{object}	PaginatedResponseModel[labelModel]	"Paginated list of labels"
//	@Failure		400			{object}	httpError							"Bad Request - Invalid query parameters"
//	@Failure		500			{object}	httpError							"Internal Server Error"
//	@Router			/labels [get]
func (e LabelGetAllEndpoint) Handle(c *gin.Context) {
	ownerTypeParams := c.QueryArray("owner_type")
	var ownerTypes []auth.OwnerType
	if len(ownerTypeParams) > 0 {
		for _, ownerTypeStr := range ownerTypeParams {
			ownerTypeStr = strings.TrimSpace(ownerTypeStr)
			if ownerType, err := auth.ParseOwnerType(ownerTypeStr); err == nil {
				ownerTypes = append(ownerTypes, ownerType)
			}
		}
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	q := query.NewFindAllQuery(ownerTypes, size, page)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[label.Label]](func(paginatedResult core.PaginatedResponse[label.Label]) any {
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

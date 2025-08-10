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
//	@Param			limit		query		integer								false	"Number of items (default: 10)"
//	@Param			offset		query		integer								false	"Offset (default: 0)"
//	@Param			familyId	query		string								false	"Family ID (UUID format)"
//	@Success		200			{object}	PaginatedResponseModel[labelModel]	"Paginated list of labels"
//	@Failure		400			{object}	HttpErrorResponse					"Bad Request - Invalid query parameters"
//	@Failure		500			{object}	HttpErrorResponse					"Internal Server Error"
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

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	if err != nil {
		offset = 0
	}

	familyIdParam := c.Query("familyId")
	familyId, err := parseUuidOrNil(&familyIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	q := query.NewFindAllQuery(ownerTypes, int32(limit), int32(offset), familyId)
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

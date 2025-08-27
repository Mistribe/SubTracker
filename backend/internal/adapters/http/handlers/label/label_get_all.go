package label

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type LabelGetAllEndpoint struct {
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[label.Label]]
}

func NewLabelGetAllEndpoint(handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[label.Label]]) *LabelGetAllEndpoint {
	return &LabelGetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all labels
//	@Description	Retrieve a paginated list of labels with optional filtering by owner type and search text
//	@Tags			labels
//	@Produce		json
//	@Param			search	query		string									false	"Search text to filter labels by name"
//	@Param			limit	query		integer									false	"Maximum number of items to return (default: 10)"
//	@Param			offset	query		integer									false	"Number of items to skip for pagination (default: 0)"
//	@Success		200		{object}	dto.PaginatedResponseModel[labelModel]	"Paginated list of labels"
//	@Failure		400		{object}	HttpErrorResponse						"Bad Request - Invalid query parameters"
//	@Failure		500		{object}	HttpErrorResponse						"Internal Server Error"
//	@Router			/labels [get]
func (e LabelGetAllEndpoint) Handle(c *gin.Context) {
	searchText := c.DefaultQuery("search", "")

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	q := query.NewFindAllQuery(searchText, limit, offset)
	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[shared.PaginatedResponse[label.Label]](func(paginatedResult shared.PaginatedResponse[label.Label]) any {
			return dto.NewPaginatedResponseModel(paginatedResult, newLabelModel)
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

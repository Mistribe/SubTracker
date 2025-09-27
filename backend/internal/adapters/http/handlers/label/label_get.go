package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type GetEndpoint struct {
	handler ports.QueryHandler[query.FindOneQuery, label.Label]
}

// Handle godoc
//
//	@Summary		Get label by LabelID
//	@Description	Retrieve a single label by its unique identifier
//	@Tags			labels
//	@Produce		json
//	@Param			id	path		string	true	"Label LabelID (UUID format)"
//	@Success		200	{object}	dto.LabelModel
//	@Failure		400	{object}	HttpErrorResponse	"Bad Request - Invalid LabelID format"
//	@Failure		404	{object}	HttpErrorResponse	"Label not found"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/{id} [get]
func (s GetEndpoint) Handle(c *gin.Context) {
	id, err := QueryParamAsUUID(c, "id")
	if err != nil {
		FromError(c, err)
		return
	}
	q := query.NewFindOneQuery(id)
	r := s.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[label.Label](func(lab label.Label) any {
			return dto.NewLabelModel(lab)
		}))
}

func (s GetEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (s GetEndpoint) Method() string {
	return http.MethodGet
}

func (s GetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewGetEndpoint(handler ports.QueryHandler[query.FindOneQuery, label.Label]) *GetEndpoint {
	return &GetEndpoint{
		handler: handler,
	}
}

package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type LabelGetEndpoint struct {
	handler ports.QueryHandler[query.FindOneQuery, label.Label]
}

// Handle godoc
//
//	@Summary		Get label by ID
//	@Description	Retrieve a single label by its unique identifier
//	@Tags			label
//	@Produce		json
//	@Param			id	path		string	true	"Label ID (UUID format)"
//	@Success		200	{object}	labelModel
//	@Failure		400	{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		404	{object}	HttpErrorResponse	"Label not found"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/{id} [get]
func (s LabelGetEndpoint) Handle(c *gin.Context) {
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
			return newLabelModel(lab)
		}))
}

func (s LabelGetEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (s LabelGetEndpoint) Method() string {
	return http.MethodGet
}

func (s LabelGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewLabelGetEndpoint(handler ports.QueryHandler[query.FindOneQuery, label.Label]) *LabelGetEndpoint {
	return &LabelGetEndpoint{
		handler: handler,
	}
}

package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

type DefaultLabelEndpoint struct {
	handler ports.QueryHandler[query.DefaultLabelQuery, []label.Label]
}

func NewDefaultLabelEndpoint(handler ports.QueryHandler[query.DefaultLabelQuery, []label.Label]) *DefaultLabelEndpoint {
	return &DefaultLabelEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get default labels
//	@Description	Retrieves a list of default system labels available to all users
//	@Tags			labels
//	@Produce		json
//	@Success		200	{array}		labelModel			"List of default labels"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/default [get]
func (e DefaultLabelEndpoint) Handle(c *gin.Context) {
	r := e.handler.Handle(c, query.DefaultLabelQuery{})
	FromResult(c,
		r,
		WithMapping[[]label.Label](func(lbls []label.Label) any {
			return slicesx.Select[label.Label, labelModel](
				lbls,
				func(lbl label.Label) labelModel {
					return newLabelModel(lbl)
				},
			)
		}),
	)
}

func (e DefaultLabelEndpoint) Pattern() []string {
	return []string{
		"/default",
	}
}

func (e DefaultLabelEndpoint) Method() string {
	return http.MethodGet
}

func (e DefaultLabelEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

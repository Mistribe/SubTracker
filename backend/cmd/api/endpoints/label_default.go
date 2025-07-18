package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/ext"
)

type DefaultLabelEndpoint struct {
	handler core.QueryHandler[query.DefaultLabelQuery, []label.Label]
}

func NewDefaultLabelEndpoint(handler core.QueryHandler[query.DefaultLabelQuery, []label.Label]) *DefaultLabelEndpoint {
	return &DefaultLabelEndpoint{handler: handler}
}

// Handle godoc
// @Summary Get default labels
// @Description Retrieves a list of default labels
// @Tags label
// @Produce json
// @Success 200 {array} labelModel
// @Router /default [get]
func (e DefaultLabelEndpoint) Handle(c *gin.Context) {
	r := e.handler.Handle(c, query.DefaultLabelQuery{})
	handleResponse(c,
		r,
		withMapping[[]label.Label](func(lbls []label.Label) any {
			return ext.Map[label.Label, labelModel](
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

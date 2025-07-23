package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/ext"
)

type LabelGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, []label.Label]
}

func NewLabelGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, []label.Label]) *LabelGetAllEndpoint {
	return &LabelGetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all labels
//	@Description	Get all labels
//	@Tags			label
//	@Produce		json
//	@Param			with_default	query		boolean	false	"Include default labels"
//	@Success		200				{array}		labelModel
//	@Failure		400				{object}	httpError
//	@Router			/labels [get]
func (e LabelGetAllEndpoint) Handle(c *gin.Context) {
	withDefault := c.DefaultQuery("with_default", "false") == "true"
	q := query.NewFindAllQuery(withDefault)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[[]label.Label](func(lbls []label.Label) any {
			return ext.Map[label.Label, labelModel](
				lbls,
				func(lbl label.Label) labelModel {
					return newLabelModel(lbl)
				},
			)
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

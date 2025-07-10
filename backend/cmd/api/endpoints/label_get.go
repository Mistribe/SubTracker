package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, label.Label]
}

func (s LabelGetEndpoint) Handle(c *gin.Context) {
	id, err := paramAsUuid(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}
	q := query.NewFindOneQuery(id)
	r := s.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[label.Label](func(lab label.Label) any {
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

func NewLabelGetEndpoint() *LabelGetEndpoint {
	return &LabelGetEndpoint{}
}

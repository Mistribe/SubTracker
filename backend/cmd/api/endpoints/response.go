package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type handleResponseOptions[TValue any] struct {
	Mapper func(TValue) any
}

type handleResponseOptionFunc[TValue any] func(*handleResponseOptions[TValue])

func withMapping[TValue any](f func(TValue) any) handleResponseOptionFunc[TValue] {
	return func(opt *handleResponseOptions[TValue]) {
		opt.Mapper = f
	}
}

func handleErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
}

func handleResponse[TValue any](c *gin.Context,
	r result.Result[TValue],
	f ...handleResponseOptionFunc[TValue]) {
	opt := handleResponseOptions[TValue]{}
	for _, v := range f {
		v(&opt)
	}

	result.Match[TValue, any](r,
		func(value TValue) any {
			if opt.Mapper != nil {
				c.JSON(http.StatusOK, opt.Mapper(value))
			} else {
				c.JSON(http.StatusOK, value)
			}
			return nil
		}, func(err error) any {
			handleErrorResponse(c, err)
			return nil
		})
}

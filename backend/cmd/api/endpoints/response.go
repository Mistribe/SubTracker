package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type handleResponseOptions[TValue any] struct {
	Mapper func(TValue) any
	Status int
}

type handleResponseOptionFunc[TValue any] func(*handleResponseOptions[TValue])

func withMapping[TValue any](f func(TValue) any) handleResponseOptionFunc[TValue] {
	return func(opt *handleResponseOptions[TValue]) {
		opt.Mapper = f
		if opt.Status == 0 {
			opt.Status = http.StatusOK
		}
	}
}

func withStatus[TValue any](status int) handleResponseOptionFunc[TValue] {
	return func(opt *handleResponseOptions[TValue]) {
		opt.Status = status
	}
}

func withNoContent[TValue any]() handleResponseOptionFunc[TValue] {
	return func(opt *handleResponseOptions[TValue]) {
		opt.Mapper = func(value TValue) any {
			return nil
		}
		opt.Status = http.StatusNoContent
	}
}

func handleErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
}

func handleResponse[TValue any](
	c *gin.Context,
	r result.Result[TValue],
	f ...handleResponseOptionFunc[TValue]) {
	opt := handleResponseOptions[TValue]{}
	for _, v := range f {
		v(&opt)
	}

	result.Match[TValue, any](r,
		func(value TValue) any {
			if opt.Mapper != nil {
				c.JSON(opt.Status, opt.Mapper(value))
			} else {
				if opt.Status == 0 {
					c.JSON(http.StatusOK, value)
				} else {
					c.JSON(opt.Status, value)
				}
			}
			return nil
		}, func(err error) any {
			handleErrorResponse(c, err)
			return nil
		})
}

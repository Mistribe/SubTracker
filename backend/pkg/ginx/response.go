package ginx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type HandleResponseOptions[TValue any] struct {
	Mapper func(TValue) any
	Status int
}

type HandleResponseOptionFunc[TValue any] func(*HandleResponseOptions[TValue])

func WithMapping[TValue any](f func(TValue) any) HandleResponseOptionFunc[TValue] {
	return func(opt *HandleResponseOptions[TValue]) {
		opt.Mapper = f
		if opt.Status == 0 {
			opt.Status = http.StatusOK
		}
	}
}

func WithStatus[TValue any](status int) HandleResponseOptionFunc[TValue] {
	return func(opt *HandleResponseOptions[TValue]) {
		opt.Status = status
	}
}

func WithNoContent[TValue any]() HandleResponseOptionFunc[TValue] {
	return func(opt *HandleResponseOptions[TValue]) {
		opt.Mapper = func(value TValue) any {
			return nil
		}
		opt.Status = http.StatusNoContent
	}
}

func FromError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, auth.ErrUnauthorized) {
		status = http.StatusUnauthorized
	}
	c.AbortWithStatusJSON(status, HttpErrorResponse{Message: err.Error()})
}

func FromResult[TValue any](
	c *gin.Context,
	r result.Result[TValue],
	f ...HandleResponseOptionFunc[TValue]) {
	opt := HandleResponseOptions[TValue]{}
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
			FromError(c, err)
			return nil
		})
}

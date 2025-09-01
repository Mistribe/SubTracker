package ginx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x/exception"
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

func httpStatusFromException(code exception.Code) int {
	switch code {
	case exception.Unknown:
		return http.StatusInternalServerError
	case exception.NotFound:
		return http.StatusNotFound
	case exception.AlreadyExists:
		return http.StatusConflict
	case exception.InvalidValue:
		return http.StatusBadRequest
	case exception.InvalidOperation:
		return http.StatusBadRequest
	case exception.InvalidState:
		return http.StatusBadRequest
	case exception.NotImplemented:
		return http.StatusNotImplemented
	case exception.InternalError:
		return http.StatusInternalServerError
	case exception.Unauthorized:
		return http.StatusForbidden
	case exception.Unauthenticated:
		return http.StatusUnauthorized
	case exception.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case exception.Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func errorResponseFromException(ex exception.Exception) (HttpErrorResponse, int) {
	status := httpStatusFromException(ex.Code())

	return HttpErrorResponse{
		Message: ex.Error(),
	}, status
}

func FromError(c *gin.Context, err error) {
	var ex exception.Exception
	if ok := errors.As(err, &ex); ok {
		response, status := errorResponseFromException(ex)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, HttpErrorResponse{Message: err.Error()})
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

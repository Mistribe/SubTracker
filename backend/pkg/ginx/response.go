package ginx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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
		opt.Mapper = func(value TValue) any { return nil }
		opt.Status = http.StatusNoContent
	}
}

// httpStatusFromException maps domain exception codes to HTTP status codes.
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

// buildProblem constructs a Problem Details object (RFC7807) for the given status & detail.
func buildProblem(c *gin.Context, status int, detail string) HttpErrorResponse {
	instance := c.FullPath()
	if instance == "" && c.Request != nil && c.Request.URL != nil {
		instance = c.Request.URL.Path
	}
	return NewProblem(ProblemTypeAboutBlank, http.StatusText(status), status, detail, instance)
}

// errorResponseFromException converts a domain exception into a Problem Details response.
func errorResponseFromException(c *gin.Context, ex exception.Exception) (HttpErrorResponse, int) {
	status := httpStatusFromException(ex.Code())
	return buildProblem(c, status, ex.Error()), status
}

// FromError writes a Problem Details response based on a generic error or domain exception.
func FromError(c *gin.Context, err error) {
	var ex exception.Exception
	if ok := errors.As(err, &ex); ok {
		response, status := errorResponseFromException(c, ex)
		c.AbortWithStatusJSON(status, response)
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		problems := make([]string, 0, len(ve))
		for _, fe := range ve {
			problems = append(problems, fmt.Sprintf("%s: %s", fe.Field(), fe.ActualTag()))
		}
		detail := strings.Join(problems, ", ")
		c.AbortWithStatusJSON(http.StatusBadRequest, buildProblem(c, http.StatusBadRequest, detail))
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, buildProblem(c, http.StatusInternalServerError, err.Error()))
}

// FromResult writes either a success response (mapped if provided) or a Problem Details error.
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

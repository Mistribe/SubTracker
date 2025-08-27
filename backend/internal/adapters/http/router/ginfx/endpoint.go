package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Endpoint interface {
	Handle(*gin.Context)
	Pattern() []string
	Method() string
	Middlewares() []gin.HandlerFunc
}

func AsEndpoint(f any) any {
	return fx.Annotate(f,
		fx.As(new(Endpoint)),
		fx.ResultTags(`group:"routes"`),
	)
}

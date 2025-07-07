package ginfx

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/fx"
)

type Route interface {
    Handle(*gin.Context)
    Pattern() []string
    Method() string
    Middlewares() []gin.HandlerFunc
}

func AsRoute(f any) any {
    return fx.Annotate(f,
        fx.As(new(Route)),
        fx.ResultTags(`group:"routes"`),
    )
}

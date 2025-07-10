package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteGroup interface {
	Prefix() string
	Routes() []Route
	Middlewares() []gin.HandlerFunc
}

func AsRouteGroup(f any) any {
	return fx.Annotate(f,
		fx.As(new(RouteGroup)),
		fx.ResultTags(`group:"route_groups"`),
	)
}

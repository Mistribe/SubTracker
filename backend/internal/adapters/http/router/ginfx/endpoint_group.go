package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type EndpointGroup interface {
	Prefix() string
	Routes() []Endpoint
	Middlewares() []gin.HandlerFunc
}

func AsEndpointGroup(f any) any {
	return fx.Annotate(f,
		fx.As(new(EndpointGroup)),
		fx.ResultTags(`group:"route_groups"`),
	)
}

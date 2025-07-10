package api

import (
	"context"
	"log/slog"
	"net/http"

	cfg "github.com/Oleexo/config-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/infrastructure/api/ginfx"
)

const (
	DefaultListenAddress = ":8080"
	ListenAddressKey     = "LISTEN_ADDRESS"
)

type EchoServerParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Logger      *slog.Logger
	Routes      []ginfx.Route      `group:"routes"`
	RouteGroups []ginfx.RouteGroup `group:"route_groups"`
	Config      cfg.Configuration
}

func registerRouteGroups(e *gin.Engine, routeGroups []ginfx.RouteGroup) {
	for _, group := range routeGroups {
		routeGroup := e.Group(group.Prefix())
		for _, m := range group.Middlewares() {
			routeGroup.Use(m)
		}
		registerRoutes(routeGroup, group.Routes())
	}
}

func registerRoutes(e *gin.RouterGroup, routes []ginfx.Route) {
	for _, route := range routes {
		for _, pattern := range route.Pattern() {
			var handlers []gin.HandlerFunc
			idx := 0
			handlers = make([]gin.HandlerFunc, len(route.Middlewares())+1)
			if route.Middlewares() != nil && len(route.Middlewares()) > 0 {
				for i, m := range route.Middlewares() {
					handlers[i] = m
					idx += 1
				}
			}
			handlers[idx] = route.Handle
			e.Handle(route.Method(), pattern, handlers...)
		}
	}
}

func newGinEngine(parameters EchoServerParams) *gin.Engine {
	e := gin.Default()

	registerRouteGroups(e, parameters.RouteGroups)
	registerRoutes(&e.RouterGroup, parameters.Routes)

	listenUrl := parameters.Config.GetStringOrDefault(ListenAddressKey, DefaultListenAddress)
	if listenUrl == "" {
		listenUrl = DefaultListenAddress
	}

	return e
}

func newHttpServer(router *gin.Engine,
	lifecycle fx.Lifecycle) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					// Handle the error (e.g., log it)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func BuildHttpServerModule() fx.Option {
	return fx.Module("ginfx",
		fx.Provide(
			newGinEngine,
			newHttpServer,
		),
		fx.Invoke(func(e *http.Server) {}),
	)
}

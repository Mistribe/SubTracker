package router

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	cfg "github.com/Oleexo/config-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	fx2 "github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

const (
	DefaultListenAddress = ":8080"
	ListenAddressKey     = "LISTEN_ADDRESS"
)

type EchoServerParams struct {
	fx.In

	Lifecycle          fx.Lifecycle
	Logger             *slog.Logger
	Routes             []fx2.Endpoint      `group:"routes"`
	RouteGroups        []fx2.EndpointGroup `group:"route_groups"`
	Config             cfg.Configuration
	LanguageMiddleware *middlewares.LanguageMiddleware
	CacheMiddleware    *middlewares.CacheMiddleware
}

func registerRouteGroups(e *gin.Engine, routeGroups []fx2.EndpointGroup) {
	for _, group := range routeGroups {
		routePrefix := group.Prefix()
		routeGroup := e.Group(routePrefix)
		for _, m := range group.Middlewares() {
			routeGroup.Use(m)
		}
		registerRoutes(routeGroup, group.Routes())
	}
}

func registerRoutes(e *gin.RouterGroup, routes []fx2.Endpoint) {
	// Order routes: static before param, param before wildcard
	type routeEntry struct {
		method      string
		pattern     string
		middlewares []gin.HandlerFunc
		handle      gin.HandlerFunc
	}

	var entries []routeEntry
	for _, route := range routes {
		for _, pattern := range route.Pattern() {
			entries = append(entries, routeEntry{
				method:      route.Method(),
				pattern:     pattern,
				middlewares: route.Middlewares(),
				handle:      route.Handle,
			})
		}
	}

	score := func(p string) int {
		switch {
		case strings.Contains(p, "/*"):
			return 2 // wildcard lowest
		case strings.Contains(p, "/:"):
			return 1 // param middle
		default:
			return 0 // static highest
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		si, sj := score(entries[i].pattern), score(entries[j].pattern)
		if si != sj {
			return si < sj
		}
		// Within the same class, register more specific (longer) paths first
		if len(entries[i].pattern) != len(entries[j].pattern) {
			return len(entries[i].pattern) > len(entries[j].pattern)
		}
		// Stable fallback by method+pattern to avoid nondeterminism
		if entries[i].method != entries[j].method {
			return entries[i].method < entries[j].method
		}
		return entries[i].pattern < entries[j].pattern
	})

	for _, it := range entries {
		var handlers []gin.HandlerFunc
		idx := 0
		handlers = make([]gin.HandlerFunc, len(it.middlewares)+1)
		if it.middlewares != nil && len(it.middlewares) > 0 {
			for i, m := range it.middlewares {
				handlers[i] = m
				idx += 1
			}
		}
		handlers[idx] = it.handle
		e.Handle(it.method, it.pattern, handlers...)
	}
}

func getOriginsFromConfig(cfg cfg.Configuration) []string {
	o := cfg.GetStringOrDefault("CORS_ALLOWED_ORIGINS", "")
	var origins []string
	if o != "" {
		origins = strings.Split(o, ",")
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:5173"}
	}
	return origins
}

func newGinEngine(parameters EchoServerParams) *gin.Engine {
	e := gin.Default()

	origins := getOriginsFromConfig(parameters.Config)

	e.Use(cors.New(cors.Config{
		AllowOrigins:     origins, // adjust as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add cache middleware
	e.Use(parameters.CacheMiddleware.Middleware())
	// Add language middleware
	e.Use(parameters.LanguageMiddleware.Middleware())

	registerRouteGroups(e, parameters.RouteGroups)
	registerRoutes(&e.RouterGroup, parameters.Routes)

	listenUrl := parameters.Config.GetStringOrDefault(ListenAddressKey, DefaultListenAddress)
	if listenUrl == "" {
		listenUrl = DefaultListenAddress
	}

	return e
}

func newHttpServer(
	router *gin.Engine,
	lifecycle fx.Lifecycle) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					// Middleware the error (e.g., log it)
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

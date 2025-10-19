package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type VersionEndpoint struct {
	version string
}

type versionParams struct {
	fx.In

	Version string `name:"appVersion"`
}

func NewVersionEndpoint(p versionParams) *VersionEndpoint {
	return &VersionEndpoint{version: p.Version}
}

// Handle godoc
//
//	@Summary		Get API version
//	@Description	Returns the build version of the SubTracker API
//	@Tags			version
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Version info"
//	@Router			/version [get]
func (e VersionEndpoint) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": e.version,
	})
}

func (e VersionEndpoint) Pattern() []string {
	return []string{
		"/version",
	}
}

func (e VersionEndpoint) Method() string {
	return http.MethodGet
}

func (e VersionEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

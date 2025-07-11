package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckLiveEndpoint struct{}

func NewHealthCheckLiveEndpoint() *HealthCheckLiveEndpoint {
	return &HealthCheckLiveEndpoint{}
}

func (e HealthCheckLiveEndpoint) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "HEALTHY",
	})
}

func (e HealthCheckLiveEndpoint) Pattern() []string {
	return []string{
		"/healthz/live",
	}
}

func (e HealthCheckLiveEndpoint) Method() string {
	return http.MethodGet
}

func (e HealthCheckLiveEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

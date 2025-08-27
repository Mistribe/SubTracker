package ginx

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func QueryParamAsUUID(c *gin.Context, key string) (uuid.UUID, error) {
	param := c.Param(key)
	if param == "" {
		return uuid.UUID{}, fmt.Errorf("missing parameter %s from route", key)
	}

	return uuid.Parse(param)

}

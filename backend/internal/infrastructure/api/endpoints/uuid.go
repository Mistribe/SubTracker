package endpoints

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func paramAsUuid(c *gin.Context, key string) (uuid.UUID, error) {
	param := c.Param(key)
	if param == "" {
		return uuid.UUID{}, fmt.Errorf("missing parameter %s from route", key)
	}

	return uuid.Parse(param)
}

func parseUuidOrNew(in *string) (uuid.UUID, error) {
	var out uuid.UUID
	var err error
	if in == nil {
		out, err = uuid.NewV7()
		if err != nil {
			return uuid.UUID{}, err
		}
	} else {
		out, err = uuid.Parse(*in)
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	return out, nil
}

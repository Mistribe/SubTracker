package provider

import (
	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Plan struct {
	*entity.Base

	name        string
	description string
	prices      []Price
}

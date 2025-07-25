package provider

import (
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Provider struct {
	*entity.Base

	name        string
	description string
	iconUrl     string
	url         string
	pricingUrl  string
	labels      []uuid.UUID
}

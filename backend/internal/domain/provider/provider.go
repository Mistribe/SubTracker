package provider

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Provider struct {
	*entity.Base

	name        string
	description string
	Labels      []uuid.UUID
}

type Plan struct {
	*entity.Base

	name        string
	description string
	Prices      []Price
}

type Price struct {
	*entity.Base

	startDate time.Time
	endDate   *time.Time
	currency  currency.Unit
	price     float64
}

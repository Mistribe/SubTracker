package provider

import (
	"time"

	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Price struct {
	*entity.Base

	startDate time.Time
	endDate   *time.Time
	currency  currency.Unit
	price     float64
}

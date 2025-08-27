package currency_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	xcur "golang.org/x/text/currency"

	dc "github.com/mistribe/subtracker/internal/domain/currency"
)

func TestNewAmount(t *testing.T) {
	a := dc.NewAmount(123.45, xcur.USD)

	assert.True(t, a.IsValid(), "new amount should be valid")
	assert.Equal(t, 123.45, a.Value(), "value should match the provided one")
	assert.Equal(t, xcur.USD, a.Currency(), "currency should match the provided unit")
}

func TestNewInvalidAmount(t *testing.T) {
	a := dc.NewInvalidAmount()

	assert.False(t, a.IsValid(), "invalid amount should not be valid")
}

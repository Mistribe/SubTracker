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

func TestAmount_IsZero(t *testing.T) {
	assert.True(t, dc.NewAmount(0, xcur.USD).IsZero())
	assert.False(t, dc.NewAmount(0.1, xcur.USD).IsZero())
}

func TestAmount_IsNegative(t *testing.T) {
	assert.True(t, dc.NewAmount(-1, xcur.USD).IsNegative())
	assert.False(t, dc.NewAmount(1, xcur.USD).IsNegative())
}

func TestAmount_IsPositive(t *testing.T) {
	assert.True(t, dc.NewAmount(1, xcur.USD).IsPositive())
	assert.False(t, dc.NewAmount(-1, xcur.USD).IsPositive())
}

func TestAmount_IsEqual(t *testing.T) {
	a1 := dc.NewAmount(100, xcur.USD)
	a2 := dc.NewAmount(100, xcur.USD)
	a3 := dc.NewAmount(200, xcur.USD)

	assert.True(t, a1.IsEqual(a2))
	assert.False(t, a1.IsEqual(a3))
}

func TestAmount_Comparisons(t *testing.T) {
	a1 := dc.NewAmount(100, xcur.USD)
	a2 := dc.NewAmount(200, xcur.USD)

	assert.True(t, a2.IsGreaterThan(a1))
	assert.True(t, a1.IsLessThan(a2))
}

func TestAmount_Arithmetic(t *testing.T) {
	a := dc.NewAmount(100, xcur.USD)

	assert.Equal(t, 150.0, a.Add(50).Value())
	assert.Equal(t, 50.0, a.Sub(50).Value())
	assert.Equal(t, 200.0, a.Mul(2).Value())
	assert.Equal(t, 50.0, a.Div(2).Value())
}

func TestAmount_WithSource(t *testing.T) {
	source := dc.NewAmount(100, xcur.USD)
	a := dc.NewAmountWithSource(200, xcur.USD, source)

	assert.Equal(t, source, a.Source())
}

func TestAmount_InvalidArithmetic(t *testing.T) {
	a := dc.NewInvalidAmount()

	assert.False(t, a.Add(50).IsValid())
	assert.False(t, a.Sub(50).IsValid())
	assert.False(t, a.Mul(2).IsValid())
	assert.False(t, a.Div(2).IsValid())
}

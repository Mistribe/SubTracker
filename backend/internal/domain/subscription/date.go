package subscription

import (
	"time"
)

func monthsBetweenCalendar(a, b time.Time) int {
	a = a.In(time.UTC)
	b = b.In(time.UTC)
	// Ensure a <= b
	if a.After(b) {
		a, b = b, a
	}
	return (b.Year()-a.Year())*12 + int(b.Month()-a.Month())
}

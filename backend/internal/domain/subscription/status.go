package subscription

import (
	"fmt"
)

type Status int8

const (
	StatusUnknown Status = iota
	StatusActive
	StatusEnded
	StatusNotStarted
)

const (
	StatusUnknownString    = "unknown"
	StatusActiveString     = "active"
	StatusEndedString      = "ended"
	StatusNotStartedString = "not_started"
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return StatusUnknownString
	case StatusActive:
		return StatusActiveString
	case StatusEnded:
		return StatusEndedString
	case StatusNotStarted:
		return StatusNotStartedString
	default:
		return StatusUnknownString
	}
}

func ParseStatus(input string) (Status, error) {
	switch input {
	case StatusActiveString:
		return StatusActive, nil
	case StatusEndedString:
		return StatusEnded, nil
	case StatusNotStartedString:
		return StatusNotStarted, nil
	default:
		return StatusUnknown, ErrUnknownStatus
	}
}

func TryParseStatus(input string) (Status, bool) {
	switch input {
	case StatusActiveString:
		return StatusActive, true
	case StatusEndedString:
		return StatusEnded, true
	case StatusNotStartedString:
		return StatusNotStarted, true
	default:
		return StatusUnknown, false
	}
}

func MustParseStatus(input string) Status {
	s, err := ParseStatus(input)
	if err != nil {
		panic(fmt.Errorf("failed to parse status: %w", err))
	}
	return s
}

package subscription

import (
	"errors"
	"strings"
)

type SortableField string

const (
	UnknownSortableField      SortableField = "unknown"
	FriendlyNameSortableField SortableField = "name"
	ProviderNameSortableField SortableField = "provider_name"
	StartDateSortableField    SortableField = "dates"
	RecurrencySortableField   SortableField = "recurrency"
)

var (
	ErrSortableFieldInvalid = errors.New("invalid sortable field")
)

func (s SortableField) String() string {
	return string(s)
}

func ParseSortableField(s string) (SortableField, error) {
	result := ParseSortableFieldOrDefault(s, UnknownSortableField)
	if result == UnknownSortableField {
		return UnknownSortableField, ErrSortableFieldInvalid
	}
	return result, nil
}

func ParseSortableFieldOrDefault(s string, defaultValue SortableField) SortableField {
	switch strings.ToLower(s) {
	case "name":
		return FriendlyNameSortableField
	case "provider":
		return ProviderNameSortableField
	case "dates":
		return StartDateSortableField
	case "recurrency":
		return RecurrencySortableField
	default:
		return defaultValue
	}
}

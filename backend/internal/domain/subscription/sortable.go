package subscription

import (
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

func (s SortableField) String() string {
	return string(s)
}

func ParseSortableField(s string) SortableField {
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
		return UnknownSortableField
	}
}

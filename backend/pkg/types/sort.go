package types

import (
	"errors"
	"strings"
)

type SortOrder string

const (
	SortOrderNone SortOrder = ""
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "ASC"
)

var (
	ErrSortOrderInvalid = errors.New("invalid sort order")
)

func (s SortOrder) String() string {
	return string(s)
}

func ParseSortOrder(s string) (SortOrder, error) {
	result := ParseSortOrderOrDefault(s, SortOrderNone)
	if result == SortOrderNone {
		return SortOrderNone, ErrSortOrderInvalid
	}
	return result, nil
}

func ParseSortOrderOrDefault(s string, defaultValue SortOrder) SortOrder {
	switch strings.ToUpper(s) {
	case "ASC":
		return SortOrderAsc
	case "DESC":
		return SortOrderDesc
	default:
		return defaultValue
	}
}

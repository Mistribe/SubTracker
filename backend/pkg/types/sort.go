package types

import (
	"strings"
)

type SortOrder string

const (
	SortOrderNone SortOrder = ""
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "ASC"
)

func (s SortOrder) String() string {
	return string(s)
}

func ParseSortOrder(s string) SortOrder {
	switch strings.ToUpper(s) {
	case "ASC":
		return SortOrderAsc
	case "DESC":
		return SortOrderDesc
	default:
		return SortOrderNone
	}
}

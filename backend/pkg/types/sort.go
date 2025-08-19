package types

type SortOrder string

const (
	SortOrderNone SortOrder = ""
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) String() string {
	return string(s)
}

func ParseSortOrder(s string) SortOrder {
	switch s {
	case "asc":
		return SortOrderAsc
	case "desc":
		return SortOrderDesc
	default:
		return SortOrderNone
	}
}

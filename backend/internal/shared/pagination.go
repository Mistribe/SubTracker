package shared

type PaginatedResponse[TValue any] struct {
	data  []TValue
	total int64
}

func (p PaginatedResponse[TValue]) Length() int {
	return len(p.data)
}

func (p PaginatedResponse[TValue]) Data() []TValue {
	return p.data
}

func (p PaginatedResponse[TValue]) Total() int64 {
	return p.total
}

func NewPaginatedResponse[TValue any](data []TValue, total int64) PaginatedResponse[TValue] {
	return PaginatedResponse[TValue]{data: data, total: total}
}

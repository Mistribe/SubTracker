package shared

type PaginatedResponse[TValue any] interface {
	Length() int
	Data() []TValue
	Total() int64
	Limits() Limits
}

type paginatedResponse[TValue any] struct {
	data   []TValue
	total  int64
	limits Limits
}

func (p paginatedResponse[TValue]) Length() int {
	return len(p.data)
}

func (p paginatedResponse[TValue]) Data() []TValue {
	return p.data
}

func (p paginatedResponse[TValue]) Total() int64 {
	return p.total
}

func (p paginatedResponse[TValue]) Limits() Limits {
	return p.limits
}

func NewPaginatedResponse[TValue any](
	data []TValue,
	total int64,
	limits Limits) PaginatedResponse[TValue] {
	return paginatedResponse[TValue]{data: data, total: total, limits: limits}
}

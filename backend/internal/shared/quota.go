package shared

type Limit interface {
	Category() string
	Feature() string
	HasLimit() bool
	Remaining() int64
	Limit() int64
}

type Limits []Limit

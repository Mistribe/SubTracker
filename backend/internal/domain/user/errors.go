package user

import (
	"github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrUnknownPlan  = exception.New("unknown plan")
	ErrUnauthorized = exception.NewUnauthorized("unauthorized")
)

func NewLimitExceededErr(feature Feature, limit int64, usage int64) PlanLimitExceeded {
	return planLimitExceeded{
		feature: feature,
		limit:   limit,
		usage:   usage,
	}
}

type PlanLimitExceeded interface {
	error
	Feature() Feature
	Limit() int64
	Usage() int64
}

type planLimitExceeded struct {
	feature Feature
	limit   int64
	usage   int64
}

func (p planLimitExceeded) Error() string {
	return "plan limit exceeded"
}

func (p planLimitExceeded) Feature() Feature {
	return p.feature
}

func (p planLimitExceeded) Limit() int64 {
	return p.limit
}

func (p planLimitExceeded) Usage() int64 {
	return p.usage
}

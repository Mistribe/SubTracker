package currency

import (
	"context"
)

type Service interface {
	ConvertTo(ctx context.Context, from Amount, to Unit) (Amount, error)
}

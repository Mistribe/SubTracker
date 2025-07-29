package endpoints

import (
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/query"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, provider.Provider]
}

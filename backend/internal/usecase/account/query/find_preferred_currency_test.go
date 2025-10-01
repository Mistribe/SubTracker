package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/account/query"
)

func TestFindPreferredCurrencyQueryHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns currency", func(t *testing.T) {
		uid := types.UserID("user-ok")

		acctService := ports.NewMockAccountService(t)
		auth := ports.NewMockAuthentication(t)
		connected := account.NewMockConnectedAccount(t)

		connected.EXPECT().UserID().Return(uid)
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(connected)
		acctService.EXPECT().GetPreferredCurrency(mock.Anything, uid).Return(currency.EUR)

		h := query.NewFindPreferredCurrencyQueryHandler(acctService, auth)
		res := h.Handle(ctx, query.NewFindPreferredCurrencyQuery())

		require.True(t, res.IsSuccess())
		var got currency.Unit
		res.IfSuccess(func(c currency.Unit) { got = c })
		assert.Equal(t, currency.EUR, got)
	})
}

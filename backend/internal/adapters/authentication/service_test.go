package authentication

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/types"
	"golang.org/x/text/currency"
)

// helper to build a minimal connected account
func newTestAccount() account.Account {
	return account.New(types.UserID("user-123"), currency.USD, types.PlanFree, types.RoleUser, nil)
}

func TestGetAccountFromContext(t *testing.T) {
	acc := newTestAccount()

	t.Run("returns account when present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ContextConnectedAccountKey, acc)
		got, ok := GetAccountFromContext(ctx)
		require.True(t, ok)
		require.Equal(t, acc, got)
	})

	t.Run("returns false when absent", func(t *testing.T) {
		ctx := context.Background()
		_, ok := GetAccountFromContext(ctx)
		require.False(t, ok)
	})
}

func TestAuthentication_MustGetConnectedAccount(t *testing.T) {
	service := NewAuthentication()
	acc := newTestAccount()

	t.Run("returns account when present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ContextConnectedAccountKey, acc)
		got := service.MustGetConnectedAccount(ctx)
		require.Equal(t, acc, got)
	})

	t.Run("panics when account missing", func(t *testing.T) {
		ctx := context.Background()
		require.PanicsWithValue(t, "missing account from context", func() {
			service.MustGetConnectedAccount(ctx)
		})
	})
}

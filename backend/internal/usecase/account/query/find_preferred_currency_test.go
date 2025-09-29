package query_test

import (
	"context"
	"errors"
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

// minimal concrete account to satisfy account.Account
type testAccount struct {
	userID types.UserID
	cur    currency.Unit
	planID types.PlanID
	role   types.Role
}

func (a testAccount) Id() string                  { return string(a.userID) }
func (a testAccount) UserID() types.UserID        { return a.userID }
func (a testAccount) Currency() currency.Unit     { return a.cur }
func (a testAccount) SetCurrency(c currency.Unit) { a.cur = c }
func (a testAccount) PlanID() types.PlanID        { return a.planID }
func (a testAccount) FamilyID() *types.FamilyID   { return nil }
func (a testAccount) Role() types.Role            { return a.role }

func newTestAccount(uid types.UserID, cur currency.Unit) account.Account {
	return testAccount{userID: uid, cur: cur, planID: types.PlanFree, role: types.RoleUser}
}

func TestFindPreferredCurrencyQueryHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("repository error returns failure", func(t *testing.T) {
		repo := ports.NewMockAccountRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := query.NewFindPreferredCurrencyQueryHandler(repo, auth)

		uid := types.UserID("user-err")
		acc := newTestAccount(uid, currency.USD)
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(acc)
		expectedErr := errors.New("boom")
		repo.EXPECT().GetById(mock.Anything, uid).Return(nil, expectedErr)

		res := h.Handle(ctx, query.NewFindPreferredCurrencyQuery())
		require.True(t, res.IsFaulted())
		var gotErr error
		res.IfFailure(func(e error) { gotErr = e })
		assert.ErrorIs(t, gotErr, expectedErr)
	})

	t.Run("success returns currency", func(t *testing.T) {
		repo := ports.NewMockAccountRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := query.NewFindPreferredCurrencyQueryHandler(repo, auth)

		uid := types.UserID("user-ok")
		acc := newTestAccount(uid, currency.EUR)
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(acc)
		repo.EXPECT().GetById(mock.Anything, uid).Return(acc, nil)

		res := h.Handle(ctx, query.NewFindPreferredCurrencyQuery())
		require.True(t, res.IsSuccess())
		var got currency.Unit
		res.IfSuccess(func(c currency.Unit) { got = c })
		assert.Equal(t, currency.EUR, got)
	})
}

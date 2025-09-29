package query_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
)

func buildSingleLabel(name string, owner types.Owner) label.Label {
	id := types.LabelID(uuid.Must(uuid.NewV7()))
	return label.NewLabel(id, owner, name, nil, "#FFFFFF", time.Now(), time.Now())
}

func TestLabelFindOneQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	t.Run("returns fault when repository errors", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.LabelID(uuid.Must(uuid.NewV7()))
		repo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, errors.New("db"))

		h := query.NewFindOneQueryHandler(repo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when label not found", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.LabelID(uuid.Must(uuid.NewV7()))
		repo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, nil)

		h := query.NewFindOneQueryHandler(repo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when label found", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		lbl := buildSingleLabel("Label A", types.NewPersonalOwner(userID))
		repo.EXPECT().GetByIdForUser(t.Context(), userID, lbl.Id()).Return(lbl, nil)

		h := query.NewFindOneQueryHandler(repo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(lbl.Id()))
		assert.True(t, res.IsSuccess())
		var got label.Label
		res.IfSuccess(func(l label.Label) { got = l })
		assert.Equal(t, lbl.Id(), got.Id())
	})
}

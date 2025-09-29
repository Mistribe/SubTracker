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
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
)

func buildLabel(name string, owner types.Owner) label.Label {
	id := types.LabelID(uuid.Must(uuid.NewV7()))
	return label.NewLabel(id, owner, name, nil, "#FFFFFF", time.Now(), time.Now())
}

func TestLabelFindAllQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	buildQuery := func() query.FindAllQuery { return query.NewFindAllQuery("lab", 10, 0) }

	t.Run("returns fault when repository errors", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		repo.EXPECT().GetAll(t.Context(), userID, ports.NewLabelQueryParameters("lab", 10, 0)).Return(nil, int64(0), errors.New("db"))

		h := query.NewFindAllQueryHandler(repo, auth, authz)
		res := h.Handle(t.Context(), buildQuery())
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success with pagination", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		labels := []label.Label{buildLabel("Label A", types.NewPersonalOwner(userID)), buildLabel("Label B", types.NewPersonalOwner(userID))}
		repo.EXPECT().GetAll(t.Context(), userID, ports.NewLabelQueryParameters("lab", 10, 0)).Return(labels, int64(2), nil)

		h := query.NewFindAllQueryHandler(repo, auth, authz)
		res := h.Handle(t.Context(), buildQuery())
		assert.True(t, res.IsSuccess())
		var pg shared.PaginatedResponse[label.Label]
		res.IfSuccess(func(v shared.PaginatedResponse[label.Label]) { pg = v })
		assert.Equal(t, int64(2), pg.Total())
		assert.Len(t, pg.Data(), 2)
	})
}

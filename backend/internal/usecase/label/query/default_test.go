package query_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
)

func buildSystemLabel(name string) label.Label {
	id := types.LabelID(uuid.Must(uuid.NewV7()))
	return label.NewLabel(id, nil, name, nil, "#000000", time.Now(), time.Now())
}

func TestDefaultLabelQueryHandler_Handle(t *testing.T) {
	// repository error case
	// success with labels
	// failure when missing defaults

	t.Run("returns fault when repository errors", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		repo.EXPECT().GetSystemLabels(t.Context()).Return(nil, errors.New("db"))
		h := query.NewDefaultLabelQueryHandler(repo)
		res := h.Handle(t.Context(), query.DefaultLabelQuery{})
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when system labels exist", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		labels := []label.Label{buildSystemLabel("A"), buildSystemLabel("B")}
		repo.EXPECT().GetSystemLabels(t.Context()).Return(labels, nil)
		h := query.NewDefaultLabelQueryHandler(repo)
		res := h.Handle(t.Context(), query.DefaultLabelQuery{})
		assert.True(t, res.IsSuccess())
		var got []label.Label
		res.IfSuccess(func(l []label.Label) { got = l })
		assert.Len(t, got, 2)
	})

	t.Run("returns fault when no system labels", func(t *testing.T) {
		repo := ports.NewMockLabelRepository(t)
		repo.EXPECT().GetSystemLabels(t.Context()).Return([]label.Label{}, nil)
		h := query.NewDefaultLabelQueryHandler(repo)
		res := h.Handle(t.Context(), query.DefaultLabelQuery{})
		assert.True(t, res.IsFaulted())
	})
}

package persistence

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelTask struct {
	labelRepository label.Repository
}

func (l LabelTask) Priority() int {
	return 10000
}

func newLabelTask(labelRepository label.Repository) *LabelTask {
	return &LabelTask{labelRepository: labelRepository}
}

func (l LabelTask) OnStart(ctx context.Context) error {
	// todo get all with limit 1 (no need to get all the database)
	lbls, err := l.labelRepository.GetAll(ctx, true)
	if err != nil {
		return err
	}
	if len(lbls) > 0 {
		return nil
	}
	for _, lbl := range getDefaultLabels() {
		if err := l.labelRepository.Save(ctx, &lbl); err != nil {
			return err
		}
	}
	return nil
}

func (l LabelTask) OnStop(_ context.Context) error {
	return nil
}

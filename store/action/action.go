package action

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/store"
)

func New(db *gorm.DB) core.ActionStore {
	return &actionStore{db: db}
}

type actionStore struct {
	db *gorm.DB
}

func (a *actionStore) ListOrderActions(ctx context.Context, checkpoint, limit int) (actions []*core.Action, err error) {
	if err = a.db.Model(core.Action{}).
		Where("id > ?", checkpoint).
		Order("id").
		Limit(limit).Find(&actions).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return actions, nil
}

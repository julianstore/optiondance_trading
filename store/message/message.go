package message

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
)

func New(db *gorm.DB) core.MessageStore {
	return &messageStore{db: db}
}

type messageStore struct {
	db *gorm.DB
}

func (s *messageStore) Create(ctx context.Context, message *core.Message) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(message).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *messageStore) CreateBatch(ctx context.Context, messages []*core.Message) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, msg := range messages {
			if err := tx.Create(msg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *messageStore) List(ctx context.Context, limit int) ([]*core.Message, error) {
	var messages []*core.Message
	if err := s.db.Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageStore) Delete(ctx context.Context, messages []*core.Message) error {
	ids := make([]int64, len(messages))
	for idx, msg := range messages {
		ids[idx] = msg.ID
	}

	if len(ids) == 0 {
		return nil
	}

	return s.db.Where("id IN (?)", ids).Delete(core.Message{}).Error
}

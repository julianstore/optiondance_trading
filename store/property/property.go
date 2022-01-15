package property

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"option-dance/core"
	"time"
)

func New(db *gorm.DB) core.PropertyStore {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (s *store) ReadProperty(ctx context.Context, key string) (string, error) {
	var p core.Property
	if err := s.db.Model(core.Property{}).Where("`key`=?", key).First(&p).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return p.Value, err
	}
	return p.Value, nil
}

func (s *store) WriteProperty(ctx context.Context, key, value string) error {
	d := core.Property{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	var count int64
	tx := s.db.Model(core.Property{}).Where("`key`=?", key).Count(&count)
	if err := tx.Error; err != nil {
		return err
	}
	if count > 0 {
		tx := s.db.Model(core.Property{}).Where("`key`=?", key).Updates(d)
		if err := tx.Error; err != nil {
			return err
		}
	}
	if count == 0 {
		return s.db.Model(core.Property{}).Create(d).Error
	}
	return nil
}

func (s *store) WriteTimeProperty(ctx context.Context, key string, value time.Time) error {
	return s.WriteProperty(ctx, key, value.UTC().Format(time.RFC3339Nano))
}

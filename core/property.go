package core

import (
	"context"
	"time"
)

type (
	Property struct {
		Key       string    `gorm:"primaryKey;column:key;type:varchar(60);not null" json:"key"`
		Value     string    `gorm:"column:value;type:varchar(100)" json:"value"`
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	}

	PropertyStore interface {
		ReadProperty(ctx context.Context, key string) (string, error)
		WriteProperty(ctx context.Context, key, value string) error
		WriteTimeProperty(ctx context.Context, key string, value time.Time) error
	}
)

func (m *Property) TableName() string {
	return "property"
}

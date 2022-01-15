package core

import (
	"context"
	"time"
)

type (
	Action struct {
		ID        int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		OrderId   string    `gorm:"column:order_id;type:varchar(45)" json:"order_id"`
		Action    string    `gorm:"column:action;type:varchar(45)" json:"action"`
		OrderData string    `gorm:"column:order_data;type:mediumtext" json:"order_data"`
		CreatedAt time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		Order     Order     `gorm:"-" json:"order"`
	}
	ActionStore interface {
		ListOrderActions(ctx context.Context, checkpoint, limit int) (actions []*Action, err error)
	}
)

func (m *Action) TableName() string {
	return "action"
}

func (m *Action) GetFilterKey() string {
	if m == nil {
		return ""
	}
	return m.OrderId + ":" + m.Action
}

package core

import (
	"context"
	"time"
)

type (
	RawTransaction struct {
		ID        int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		CreatedAt time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		TraceID   string    `gorm:"column:trace_id;type:varchar(45)" json:"trace_id"`
		UserID    string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		Amount    string    `gorm:"column:amount;type:varchar(45)" json:"amount"`
		AssetID   string    `gorm:"column:asset_id;type:varchar(45)" json:"asset_id"`
		Data      string    `gorm:"column:data;type:mediumtext" json:"data"`
		State     string    `gorm:"column:state;type:varchar(45)" json:"state"`
		ErrorInfo string    `gorm:"column:error_info;type:varchar(200)" json:"error_info"`
	}
	RawTxStore interface {
		ListPendingRawTransactions(ctx context.Context, limit int) (result []*RawTransaction, err error)
		ExpireRawTransaction(ctx context.Context, tx *RawTransaction) error
		SaveRawTransaction(ctx context.Context, tx *RawTransaction) error
	}
)

func (m *RawTransaction) TableName() string {
	return "raw_transaction"
}

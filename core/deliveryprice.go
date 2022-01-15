package core

import (
	"context"
	"time"
)

type (
	DeliveryPrice struct {
		ID        int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		Asset     string    `gorm:"column:asset;type:varchar(45);" json:"asset"`
		Date      string    `gorm:"column:date;type:varchar(45);unique;" json:"date"`
		Price     string    `gorm:"column:price;type:varchar(45);" json:"price"`
		CreatedAt time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	}
	DeliveryPriceStore interface {
		WritePrice(ctx context.Context, asset, date, price string, time time.Time) error
		ReadPrice(ctx context.Context, asset, date string) (string, error)
		ListPricesByAsset(ctx context.Context, asset string) ([]*DeliveryPrice, error)
	}
)

func (m *DeliveryPrice) TableName() string {
	return "delivery_price"
}

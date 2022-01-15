package core

import (
	"context"
	"time"
)

const (
	DbBtcUsdIndexPrice = "dbBtcUsdIndexPrice"
)

type (
	DbMarket struct {
		ID                  int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		InstrumentName      string    `gorm:"column:instrument_name;type:varchar(45);not null;unique" json:"instrument_name"`
		StrikePrice         int64     `gorm:"column:strike_price;type:bigint" json:"strike_price"`
		ExpirationDate      time.Time `gorm:"column:expiration_date;type:datetime" json:"expiration_date"`
		ExpirationDateStr   string    `gorm:"column:expiration_date_str;type:varchar(45)" json:"expiration_date_str"`
		ExpirationTimestamp int64     `gorm:"column:expiration_timestamp;type:bigint" json:"expiration_timestamp"`
		OptionType          string    `gorm:"column:option_type;type:varchar(45)" json:"option_type"`
		DeliveryType        string    `gorm:"column:delivery_type;type:varchar(45)" json:"delivery_type"`
		QuoteCurrency       string    `gorm:"column:quote_currency;type:varchar(45)" json:"quote_currency"`
		BaseCurrency        string    `gorm:"column:base_currency;type:varchar(45)" json:"base_currency"`
		Bids                string    `gorm:"column:bids;type:text" json:"-"`
		Asks                string    `gorm:"column:asks;type:text" json:"-"`
		CreatedAt           time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
		Status              int8      `gorm:"column:status;type:tinyint" json:"status"`
	}

	DbMarketService interface {
		MarketService
		SyncMarket(ctx context.Context) error
		SyncIndexPrice(ctx context.Context) error
	}

	DbMarketStore interface {
		MarketStore
		Save(ctx context.Context, market DbMarket) error
		CloseMarket(ctx context.Context, name string) error
		ListAll(ctx context.Context) ([]*DbMarket, error)
	}
)

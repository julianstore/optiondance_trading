package core

import (
	"context"
	"time"
)

const (
	TrackStatusOn  = 0
	TrackStatusOff = 1
)

type (
	DbMarketTracker struct {
		ID            int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		Date          string    `gorm:"column:date;type:varchar(45);" json:"date"`
		Type          string    `gorm:"column:type;type:varchar(45);" json:"type"`
		DeliveryType  string    `gorm:"column:delivery_type;type:varchar(45)" json:"delivery_type"`
		StrikePrices  string    `gorm:"column:strike_prices;type:varchar(500);" json:"strike_prices"`
		QuoteCurrency string    `gorm:"column:quote_currency;type:varchar(45)" json:"quote_currency"`
		BaseCurrency  string    `gorm:"column:base_currency;type:varchar(45)" json:"base_currency"`
		CreatedAt     time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
		Status        int8      `gorm:"column:status;type:tinyint" json:"status"`
	}

	MarketTracker struct {
		DbInstrumentName string `json:"db_instrument_name"`
		OdInstrumentName string `json:"od_instrument_name"`
		BidTrack         bool
		AskTrack         bool
		Bids             string `gorm:"column:bids;type:text" json:"-"`
		Asks             string `gorm:"column:asks;type:text" json:"-"`
		Status           int8
	}

	DbMarketTrackerStore interface {
		Save(ctx context.Context, date, marketType, strikePrices string)
		ListAll(ctx context.Context) ([]*DbMarketTracker, error)
		ToMarketTracker(ctx context.Context, trackers []*DbMarketTracker) (marketTracker []*MarketTracker, err error)
	}
)

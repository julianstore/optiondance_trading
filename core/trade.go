package core

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"gorm.io/gorm"
	"time"
)

const (
	TradeStatusOpen   = 10
	TradeStatusClosed = 20
)

type (
	Trade struct {
		ID              int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		TradeID         string    `gorm:"column:trade_id;type:varchar(45)" json:"trade_id"`
		UserID          string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		Liquidity       string    `gorm:"column:liquidity;type:varchar(45)" json:"liquidity"`
		BidOrderID      string    `gorm:"column:bid_order_id;type:varchar(45)" json:"bid_order_id"`
		AskOrderID      string    `gorm:"column:ask_order_id;type:varchar(45)" json:"ask_order_id"`
		QuoteAssetID    string    `gorm:"column:quote_asset_id;type:varchar(45)" json:"quote_asset_id"`
		BaseAssetID     string    `gorm:"column:base_asset_id;type:varchar(45)" json:"base_asset_id"`
		Side            string    `gorm:"column:side;type:varchar(45)" json:"side"`
		Price           string    `gorm:"column:price;type:varchar(45)" json:"price"`
		Amount          string    `gorm:"column:amount;type:varchar(45)" json:"amount"`
		InstrumentName  string    `gorm:"column:instrument_name;type:varchar(45)" json:"instrument_name"`
		FeeAmount       string    `gorm:"column:fee_amount;type:varchar(45)" json:"fee_amount"`
		FeeAssetID      string    `gorm:"column:fee_asset_id;type:varchar(45)" json:"fee_asset_id"`
		UnderlyingPrice string    `gorm:"column:underlying_price;type:varchar(45)" json:"underlying_price"`
		Status          int8      `gorm:"column:status;type:tinyint" json:"status"`
		CreatedAt       time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	}

	TradeStore interface {
		ListMonthlyPremium(ctx context.Context, uid, side string, year int) (monthlyProfit []MonthlyProfit, err error)
		ListOrderTrades(userId string, orderId string, side string) (tradeList []*Trade, err error)
	}

	TradeService interface {
		Transact(ctx context.Context, taker, maker *EngineOrder, amount number.Integer) (string, error)
		CancelOrder(ctx context.Context, order *EngineOrder, cancelledAt time.Time) error
		MakeOrderMutations(taker, maker *EngineOrder, tx *gorm.DB) (err error)
		MakeTrades(taker, maker *EngineOrder, amount number.Decimal) (*Trade, *Trade)
		HandleFees(ask, bid *Trade, taker, maker *EngineOrder) (*Transfer, *Transfer)
		HandleFeesWithBidRefund(ask, bid *Trade, taker, maker *EngineOrder) (*Transfer, *Transfer, *Transfer)
	}
)

func (m *Trade) TableName() string {
	return "trade"
}

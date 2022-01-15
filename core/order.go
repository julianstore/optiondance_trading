package core

import (
	"context"
	"time"
)

const (
	OrderStatusOpen      = 10
	OrderStatusFilled    = 20
	OrderStatusRejected  = 30
	OrderStatusCanceling = 35
	OrderStatusCancelled = 40
	OrderStateDone       = 50
)

type (
	Order struct {
		ID                  int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		OrderID             string    `gorm:"column:order_id;type:varchar(45);index:idx_order_id" json:"order_id"`
		OrderNum            int64     `gorm:"column:order_num;type:bigint" json:"order_num"`
		UserID              string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		Side                string    `gorm:"column:side;type:varchar(45)" json:"side"`
		OrderType           string    `gorm:"column:order_type;type:varchar(45)" json:"order_type"`
		Price               string    `gorm:"column:price;type:varchar(45)" json:"price"`
		RemainingAmount     string    `gorm:"column:remaining_amount;type:varchar(45)" json:"remaining_amount"`
		FilledAmount        string    `gorm:"column:filled_amount;type:varchar(45)" json:"filled_amount"`
		RemainingFunds      string    `gorm:"column:remaining_funds;type:varchar(45)" json:"remaining_funds"`
		FilledFunds         string    `gorm:"column:filled_funds;type:varchar(45)" json:"filled_funds"`
		Margin              string    `gorm:"column:margin;type:varchar(45)" json:"margin"`
		QuoteAssetID        string    `gorm:"column:quote_asset_id;type:varchar(45)" json:"quote_asset_id"`
		BaseAssetID         string    `gorm:"column:base_asset_id;type:varchar(45)" json:"base_asset_id"`
		InstrumentName      string    `gorm:"column:instrument_name;type:varchar(45);index:idx_name_exp" json:"instrument_name"`
		OptionType          string    `gorm:"column:option_type;type:varchar(45)" json:"option_type"`
		DeliveryType        string    `gorm:"column:delivery_type;type:varchar(45)" json:"delivery_type"`
		StrikePrice         string    `gorm:"column:strike_price;type:varchar(45)" json:"strike_price"`
		ExpirationTimestamp int64     `gorm:"column:expiration_timestamp;type:bigint;index:idx_name_exp" json:"expiration_timestamp"`
		ExpirationDate      time.Time `gorm:"column:expiration_date;type:datetime" json:"expiration_date"`
		QuoteCurrency       string    `gorm:"column:quote_currency;type:varchar(45)" json:"quote_currency"`
		BaseCurrency        string    `gorm:"column:base_currency;type:varchar(45)" json:"base_currency"`
		OrderStatus         int8      `gorm:"column:order_status;type:tinyint" json:"order_status"`
		CreatedAt           time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
		TradeList           []*Trade  `gorm:"-" json:"trade_list"`
		OrderNumString      string    `gorm:"-" json:"order_num_string"`
	}

	OrderStore interface {
		FindByOrderId(ctx context.Context, orderId string) (orders *Order, err error)
		ListOrdersByINameAndStatus(ctx context.Context, instrumentName string, status int64) (orders []*Order, err error)
		ListOrderByIds(ctx context.Context, orderIds []string) (orders []Order, err error)
		ListInstrumentNamesByExpiryDate(ctx context.Context, time time.Time) (names []string, err error)
		ListOrdersBeforeExpiryDate(ctx context.Context, status int64, time time.Time) (openOrders []*Order, err error)
		ListOrdersBeforeTime(ctx context.Context, status int64, side string, time time.Time) (openOrders []*Order, err error)
		ListOrdersByPage(current int64, size int64, uid, status, order string) (list []*Order, total int64, pages int64)
	}

	OrderService interface {
		FindByOrderId(ctx context.Context, userId, orderId string) (orders *Order, err error)
		ListOpenOrdersByExpiryDate(ctx context.Context, time time.Time) (orders []*Order, err error)
		ListOpenOrdersBeforeExpiryDate(ctx context.Context, time time.Time) (orders []*Order, err error)
		ListOpenOrdersBeforeTime(ctx context.Context, side string, time time.Time) (openOrders []*Order, err error)
		Page(current int64, size int64, uid string, status, order string) (list []*Order, total int64, pages int64)
		ListPendingActions(ctx context.Context, checkpoint int, limit int) (actions []*Action, err error)
		CreateOrderAction(ctx context.Context, o *Order, userId, actionType string, createdAt time.Time) error
		CancelOrderAction(ctx context.Context, orderId string, createdAt time.Time) error
	}
)

func (m *Order) TableName() string {
	return "order"
}

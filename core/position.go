package core

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

const (
	PositionSideAsk  = "ASK"
	PositionSideBid  = "BID"
	PositionSideZero = "ZERO"

	PositionStatusALL               = -1
	PositionStatusNotExercised      = 10
	PositionStatusExerciseRequested = 20
	PositionStatusExercised         = 30

	ExerciseBtnStatusHidden         = "hidden"
	ExerciseBtnStatusWaitExercise   = "wait_exercise"
	ExerciseBtnStatusOpenExercise   = "open_exercise"
	ExerciseBtnStatusWaitSettlement = "wait_settlement"
)

type (
	Position struct {
		ID             int64   `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		PositionID     string  `gorm:"column:position_id;type:varchar(45);index:idx_position_id" json:"position_id"`
		PositionNum    int64   `gorm:"column:position_num;type:bigint" json:"position_num"`
		UserID         string  `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		InstrumentName string  `gorm:"column:instrument_name;type:varchar(45);index:idx_name" json:"instrument_name"`
		Side           string  `gorm:"column:side;type:varchar(45)" json:"side"`
		Type           string  `gorm:"column:type;type:varchar(45)" json:"type"`
		Size           float64 `gorm:"column:size;type:decimal(32,8)" json:"size"`
		Funds          float64 `gorm:"column:funds;type:decimal(32,8)" json:"funds"`
		Margin         float64 `gorm:"column:margin;type:decimal(32,8)" json:"margin"`
		ExercisedSize  float64 `gorm:"column:exercised_size;type:decimal(32,8)" json:"exercised_size"`
		AveragePrice   float64 `gorm:"column:average_price;type:decimal(32,8)" json:"average_price"`

		StrikePrice   string `gorm:"column:strike_price;type:varchar(45)" json:"strike_price"`
		OptionType    string `gorm:"column:option_type;type:varchar(45)" json:"option_type"`
		DeliveryType  string `gorm:"column:delivery_type;type:varchar(45)" json:"delivery_type"`
		QuoteCurrency string `gorm:"column:quote_currency;type:varchar(45)" json:"quote_currency"`
		BaseCurrency  string `gorm:"column:base_currency;type:varchar(45)" json:"base_currency"`

		ExpirationDate      time.Time `gorm:"column:expiration_date;type:datetime" json:"expiration_date"`
		ExpirationTimestamp int64     `gorm:"column:expiration_timestamp;type:bigint" json:"expiration_timestamp"`
		Status              int8      `gorm:"column:status;type:tinyint" json:"status"`
		Settlement          int8      `gorm:"column:settlement;type:tinyint;default:0" json:"settlement"`
		CreatedAt           time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
		PositionNumString   string    `gorm:"-" json:"position_num_string"`
		TradeList           []*Trade  `gorm:"-" json:"trade_list"`
		InitialFunds        float64   `gorm:"-" json:"initial_funds"`
	}

	PositionDTO struct {
		*Position
		SettlementInfo SettlementInfo `json:"settlement_info"`
		Exercisable    bool           `json:"exercisable"`
		Expiry         bool           `json:"expiry"`
		BtnStatus      string         `json:"btn_status"`
	}

	PositionStore interface {
		FindByPositionId(ctx context.Context, positionId string) (position *Position, err error)
		FindByPositionIdAndUid(ctx context.Context, uid, positionId string) (position *Position, err error)

		ListInstrumentNamesByDate(ctx context.Context, time time.Time, status int64) (names []string, err error)
		ListNamesByDateAndDeliveryType(ctx context.Context, time time.Time, status int64, deliveryType string) (names []string, err error)
		ListNamesByDateAndSide(ctx context.Context, time time.Time, status int64, side, deliveryType string) (names []string, err error)
		ListExercisablePosition(ctx context.Context, instrumentName string) (positions []*Position, err error)
		ListMonthlyPremium(ctx context.Context, uid string, year int) (monthlyProfit []MonthlyProfit, err error)
		ListMonthlyUnderlying(ctx context.Context, uid string, year int) (monthlyProfit []MonthlyProfit, err error)
		ListSidePositions(ctx context.Context, instrumentName, side string, statusSet []int64) (positions []*Position, err error)

		ListPositionTrades(uid, instrumentName, side string) (tradeList []*Trade, err error)
		FindPositionStatusById(id string) (status int8, err error)
		FindPositionByUidAndName(uid, instrumentName string) (position *Position, err error)
		ListUserPositions(current, size int64, uid, status, order string) (list []*Position, total, pages int64, err error)

		ExercisePositionWithID(ctx context.Context, positionId, instrumentName string, updateTime time.Time) (err error)
		AutoExercisePositionWithName(ctx context.Context, instrumentName string, updateTime time.Time) (err error)
		UpdateExercisedSize(id int64, size float64, updatedAt time.Time) error
	}

	PositionService interface {
		PositionDetail(ctx context.Context, uid, positionId string) (position *PositionDTO, err error)
		SettlementInfo(ctx context.Context, position *Position) (SettlementInfo, error)

		UpdatePositionWithTrade(tx *gorm.DB, trade *Trade) (err error)
		AutoExercise(ctx context.Context, dateString string, deliveryPrice decimal.Decimal, exercisedAt time.Time) error

		Settlement(ctx context.Context, asset, dateString, deliveryType string, settledAt time.Time) error
		SettlementGather(ctx context.Context, asset, dateString, deliveryType string, updateDb bool, settledAt time.Time) (GatherResult, error)
		SettlementCommit(ctx context.Context, transferList []*Transfer, exercisedNameList, nameList []string, settledAt time.Time) error
		LogGatherResult(result GatherResult)
	}
	GatherResult struct {
		NameList          []string
		ExercisedNameList []string
		Transfers         []*Transfer
	}

	SettlementInfo struct {
		Exercised       bool   `json:"exercised"`
		Underlying      string `json:"underlying"`
		UnderlyingPrice string `json:"underlying_price"`
		RefundMargin    string `json:"refund_margin"`
		BidEarnings     string `json:"bid_earnings"`
		Size            string `json:"size"`
	}
)

func (m *Position) TableName() string {
	return "position"
}

func PositionSideFromSize(size number.Decimal) (side string) {
	cmp := size.Cmp(number.Zero())
	switch {
	case cmp > 0:
		side = PositionSideBid
		break
	case cmp == 0:
		side = PositionSideZero
		break
	case cmp < 0:
		side = PositionSideAsk
		break
	}
	return
}

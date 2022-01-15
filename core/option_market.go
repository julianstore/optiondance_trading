package core

import (
	"context"
	"encoding/json"
	"time"
)

type (
	OptionMarket struct {
		InstrumentName      string    `gorm:"primaryKey;column:instrument_name;type:varchar(45);not null" json:"instrument_name"`
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

	Entry struct {
		Side   string `json:"side"`
		Price  string `json:"price"`
		Amount string `json:"amount"`
		Funds  string `json:"funds"`
	}

	OptionMarketDTO struct {
		OptionMarket
		BidList []Entry `gorm:"-" json:"bids"`
		AskList []Entry `gorm:"-" json:"asks"`
	}

	FullMarketByDateDTO struct {
		StrikePrice int             `json:"strike_price"`
		Call        OptionMarketDTO `json:"call"`
		Put         OptionMarketDTO `json:"put"`
	}

	MarketService interface {
		ListExpiryDatesByPrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []OptionMarketDTO, err error)
		ListExpiryDatesByPriceOdAndDb(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []OptionMarketDTO, err error)
		ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int64)
		ListMarketStrikePriceOdAndDb(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int)
		ListMarketsByInstrument(name string) (d OptionMarketDTO, err error)
		ListMarketsByInstrumentOdAndDb(name string) (d OptionMarketDTO, err error)
		UpdateOptionMarket(instrumentName string, bids, asks string) (err error)

		ListDates(ctx context.Context, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error)
		ListByDate(ctx context.Context, date string) (dateList []FullMarketByDateDTO, err error)
	}

	MarketStore interface {
		Create(market OptionMarket) error
		FindByInstrumentName(instrumentName string) (o OptionMarketDTO, err error)
		ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (list []OptionMarketDTO, err error)
		ListDates(ctx context.Context, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error)
		ListByDate(ctx context.Context, date string) (list []FullMarketByDateDTO, err error)
		ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarkets []OptionMarketDTO, err error)
		UpdateBidAsks(instrumentName string, bids, asks string) (rowsAffected int64, err error)
		DeleteByInstrumentName(ctx context.Context, instrumentName string) error
	}
)

func (m *OptionMarket) TableName() string {
	return "option_market"
}

func ConvertToMarketDTO(market OptionMarket) (d OptionMarketDTO, err error) {
	var bidList, askList []Entry
	err = json.Unmarshal([]byte(market.Asks), &askList)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal([]byte(market.Bids), &bidList)
	if err != nil {
		return d, err
	}
	d = OptionMarketDTO{OptionMarket: market, BidList: bidList, AskList: askList}
	return
}

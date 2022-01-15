package core

import (
	"github.com/MixinNetwork/go-number"
	"log"
	"time"
)

const (
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
	PageSideAsk     = "ASK"
	PageSideBid     = "BID"
)

type EngineOrder struct {
	Id              string
	Side            string
	Type            string
	Price           number.Integer
	RemainingAmount number.Integer
	FilledAmount    number.Integer
	RemainingFunds  number.Integer
	FilledFunds     number.Integer
	InstrumentName  string

	Quote         string
	Base          string
	QuoteCurrency string
	BaseCurrency  string
	UserId        string
	BrokerId      string
	OptionType    string
	Margin        string
	LastPrice     number.Integer
	CreatedAt     time.Time
}

func (order *EngineOrder) Filled() bool {
	//if order.Side == PageSideAsk {
	//	return order.RemainingAmount.IsZero()
	//}
	//return order.RemainingFunds.IsZero()
	return order.RemainingAmount.IsZero()
}

func (order *EngineOrder) Assert() {
	switch order.Side {
	case PageSideAsk:
		if !order.RemainingFunds.IsZero() {
			log.Panicln(order)
		}
	case PageSideBid:
		//if !order.RemainingAmount.IsZero() {
		//	log.Panicln(order)
		//}
	default:
		log.Panicln(order)
	}

	switch order.Type {
	case OrderTypeLimit:
		if order.Price.IsZero() {
			log.Panicln(order)
		}
	case OrderTypeMarket:
		if !order.Price.IsZero() {
			log.Panicln(order)
		}
	default:
		log.Panicln(order)
	}
}

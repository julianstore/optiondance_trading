package main

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"option-dance/cmd/config"
	marketz "option-dance/service/market"
	"option-dance/service/notifier"
	orderz "option-dance/service/order"
	positionz "option-dance/service/position"
	"option-dance/service/statistics"
	tradez "option-dance/service/trade"
	transferz "option-dance/service/transfer"
	userz "option-dance/service/user"
	"option-dance/service/waitlist"
	"option-dance/store/action"
	"option-dance/store/deliveryprice"
	markets "option-dance/store/market"
	orders "option-dance/store/order"
	"option-dance/store/position"
	trades "option-dance/store/trade"
	"option-dance/store/transfer"
	users "option-dance/store/user"
	waitlists "option-dance/store/waitlist"
)

func ProvideDb() *gorm.DB {
	return config.Db
}

var apiServerSet = wire.NewSet(
	ProvideDb,
	users.New,
	userz.New,
	waitlists.New,
	waitlist.New,
	deliveryprice.New,
	position.New,
	trades.New,
	transfer.New,
	transferz.New,
	notifier.NewMuteNotifier,
	tradez.NewTradeService,
	statistics.New,
	markets.New,
	marketz.New,
	positionz.New,
	action.New,
	orders.New,
	orderz.New,
)

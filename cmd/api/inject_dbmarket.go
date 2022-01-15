package main

import (
	"github.com/google/wire"
	"option-dance/service/market"
	"option-dance/store/dbmarket"
	"option-dance/store/property"
)

var apiDbMarketSet = wire.NewSet(
	property.New,
	dbmarket.NewDbMarketStore,
	dbmarket.NewDbMarketTrackerStore,
	market.NewDbMarketService,
)

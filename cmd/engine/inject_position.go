package main

import (
	"option-dance/core"
	"option-dance/service/position"
	"option-dance/service/transfer"
	"option-dance/store/deliveryprice"
	positions "option-dance/store/position"
	"option-dance/store/trade"
	transfers "option-dance/store/transfer"
	"option-dance/store/user"
	"option-dance/store/utxo"
)

func ProvidePositionService() core.PositionService {
	db := ProvideDb()
	positionStore := positions.New(db)
	tradeStore := trade.New(db)
	transferStore := transfers.New(db)
	userStore := user.New(db)
	transferService := transfer.New(transferStore)
	deliveryPriceStore := deliveryprice.New(db)
	positionService := position.New(positionStore, tradeStore, transferService, userStore, deliveryPriceStore)
	return positionService
}

func ProvideUtxoStore() core.UtxoStore {
	db := ProvideDb()
	utxoStore := utxo.New(db)
	return utxoStore
}

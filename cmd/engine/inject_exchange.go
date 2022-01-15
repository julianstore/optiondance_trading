package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/google/wire"
	"gorm.io/gorm"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/match"
	m "option-dance/pkg/mixin"
	marketz "option-dance/service/market"
	"option-dance/service/message"
	orderz "option-dance/service/order"
	positionz "option-dance/service/position"
	tradez "option-dance/service/trade"
	transferz "option-dance/service/transfer"
	userz "option-dance/service/user"
	actions "option-dance/store/action"
	"option-dance/store/dbmarket"
	"option-dance/store/deliveryprice"
	"option-dance/store/market"
	messages "option-dance/store/message"
	orders "option-dance/store/order"
	positions "option-dance/store/position"
	"option-dance/store/property"
	"option-dance/store/rawtx"
	trades "option-dance/store/trade"
	transfers "option-dance/store/transfer"
	"option-dance/store/user"
	"option-dance/store/utxo"
)

func ProvideDb() *gorm.DB {
	return config.Db
}

func ProvideUserService(userStore core.UserStore) core.UserService {
	return userz.New(userStore)
}

func ProvideMixinClient() *mixin.Client {
	client, _ := m.Client()
	return client
}

var exchangeSet = wire.NewSet(
	ProvideDb,
	ProvideMixinClient,
	transfers.New,
	transferz.New,
	user.New,
	utxo.New,
	orders.New,
	orderz.New,
	rawtx.New,
	trades.New,
	actions.New,
	positions.New,
	positionz.New,
	deliveryprice.New,
	messages.New,
	message.NewMessageBuilder,
	market.New,
	dbmarket.NewDbMarketStore,
	marketz.New,
	property.New,
	provideNotifier,
	tradez.NewTradeService,
	ProvideUserService,
	message.New,
	match.NewUtxoDispatcher,
	match.NewUtxoSyncer,
	match.NewSpentSyncer,
	match.NewMultiSigner,
)

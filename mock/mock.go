package mock

//go:generate mockgen -package=mock -destination=mock_gen.go option-dance/core PositionService,PositionStore,TradeStore

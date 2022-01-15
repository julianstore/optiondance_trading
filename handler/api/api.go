package api

import (
	"github.com/gin-gonic/gin"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/handler/api/auth"
	"option-dance/handler/api/dbmarket"
	"option-dance/handler/api/deliveryprice"
	"option-dance/handler/api/hc"
	"option-dance/handler/api/market"
	"option-dance/handler/api/order"
	"option-dance/handler/api/position"
	"option-dance/handler/api/statistics"
	"option-dance/handler/api/user"
	"option-dance/handler/api/waitlist"
	"option-dance/handler/middleware"
)

type ServerApi struct {
	UserService        core.UserService
	WaitListService    core.WaitListService
	MarketService      core.MarketService
	StatisticsService  core.StatisticsService
	orderService       core.OrderService
	positionService    core.PositionService
	positionStore      core.PositionStore
	propertyStore      core.PropertyStore
	deliveryPriceStore core.DeliveryPriceStore
	DbMarketService    core.DbMarketService
}

func NewServer(
	userService core.UserService,
	waitListService core.WaitListService,
	marketService core.MarketService,
	statisticsService core.StatisticsService,
	orderService core.OrderService,
	positionService core.PositionService,
	positionStore core.PositionStore,
	deliveryPriceStore core.DeliveryPriceStore,
	propertyStore core.PropertyStore,
	dbMarketService core.DbMarketService,
) ServerApi {
	return ServerApi{
		UserService:        userService,
		WaitListService:    waitListService,
		MarketService:      marketService,
		StatisticsService:  statisticsService,
		orderService:       orderService,
		positionStore:      positionStore,
		positionService:    positionService,
		deliveryPriceStore: deliveryPriceStore,
		DbMarketService:    dbMarketService,
		propertyStore:      propertyStore,
	}
}

func InitializeApiServer(s ServerApi) *gin.Engine {

	app := gin.Default()
	app.Use(config.GinLogger(), config.GinRecovery(true))
	app.Use(middleware.Cors())

	//auth
	app.POST("/api/v1/auth", auth.HandleAuthByMixinToken(s.UserService))
	//waitList
	app.POST("/api/v1/waitlist", waitlist.WaitListCreate(s.WaitListService))
	app.GET("/api/v1/waitlist/rank-info", waitlist.WaitlistRankInfo(s.WaitListService))
	app.GET("/api/v1/waitlist/wid", waitlist.WaitlistWid(s.WaitListService))
	//od-market
	app.GET("/api/v1/market/strike-prices", market.ListStrikePrices(s.MarketService, s.DbMarketService))
	app.GET("/api/v1/market/price/:price", market.ListExpiryDatesByPrice(s.MarketService, s.DbMarketService))
	app.GET("/api/v1/market/instrument/:name", market.FindByInstrument(s.MarketService, s.DbMarketService))

	app.GET("/api/v1/market/dates", market.ListDates(s.MarketService, s.DbMarketService))
	app.GET("/api/v1/market/date/:date", market.ListByDate(s.MarketService, s.DbMarketService))

	//db-market
	app.GET("/api/v1/db-market/strike-prices", dbmarket.ListStrikePrices(s.DbMarketService))
	app.GET("/api/v1/db-market/price/:price", dbmarket.ListExpiryDatesByPrice(s.DbMarketService))
	app.GET("/api/v1/db-market/instrument/:name", dbmarket.FindByInstrument(s.DbMarketService))
	app.GET("/api/v1/db-market/index-price", dbmarket.IndexPrice(s.propertyStore))

	//----------------------Interceptor----------------------
	app.Use(middleware.Interceptor())
	//user
	app.POST("/api/v1/user/me", user.HandlerUserMe(s.UserService))
	app.GET("/api/v1/user-settings", user.HandlerUserSettings(s.UserService))
	app.POST("/api/v1/user-settings", user.HandlerUpdateUserSettings(s.UserService))
	//order
	app.POST("/api/v1/order-request", order.OrderCreateRequest())
	app.POST("/api/v1/order/:id/cancel", order.OrderCancel(s.UserService))
	app.GET("/api/v1/order-trace/:id", order.OrderDetailByOrderId(s.UserService, s.orderService))
	app.GET("/api/v1/order/:id", order.OrderDetailByOrderId(s.UserService, s.orderService))
	app.GET("/api/v1/orders", order.Page(s.UserService, s.orderService))
	app.GET("/api/v1/open-orders", order.OpenOrders(s.orderService))
	app.GET("/api/v1/open-orders/before-now", order.OpenOrdersBeforeNow(s.orderService))

	//position
	app.POST("/api/v1/exercise-request", position.ExerciseRequest())
	app.POST("/api/v1/exercise-cash", position.ExerciseCashDelivery())
	app.GET("/api/v1/positions", position.ListUserPositions(s.UserService, s.positionStore))
	app.GET("/api/v1/position/:id", position.PositionDetail(s.UserService, s.positionStore))
	app.GET("/api/v2/position/:id", position.DetailByID(s.UserService, s.positionService))
	app.GET("/api/v1/position-name/:name", position.PositionDetailByInstrumentName(s.UserService, s.positionStore))
	app.GET("/api/v1/position-status/:id", position.PositionStatus(s.positionStore))
	//statistics analysis
	app.GET("/api/v1/statistics/annual-sell-put-premium", statistics.AnnualSellPutPremium(s.StatisticsService, s.UserService))
	app.GET("/api/v1/statistics/annual-sell-put-underlying", statistics.AnnualSellPutUnderlying(s.StatisticsService, s.UserService))

	//delivery price
	app.GET("/api/v1/delivery-price", deliveryprice.DeliveryPrice(s.deliveryPriceStore))
	app.GET("/api/v1/delivery-prices", deliveryprice.DeliveryPrices(s.deliveryPriceStore))
	return app
}

func InitEngineServer() *gin.Engine {
	app := gin.Default()
	app.GET("/hc", hc.Hc())
	return app
}

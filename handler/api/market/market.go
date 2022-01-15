package market

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/pkg/http"
)

func FindByInstrument(marketService core.MarketService, dbMarketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		instrument, err := marketService.ListMarketsByInstrumentOdAndDb(name)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(instrument, c)
	}
}

func ListStrikePrices(marketService core.MarketService, dbMarketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		side := c.Query("side")
		optionType := c.Query("optionType")
		deliveryType := c.Query("deliveryType")
		quoteCurrency := c.Query("quoteCurrency")
		baseCurrency := c.Query("baseCurrency")
		if side == "" || optionType == "" || deliveryType == "" || quoteCurrency == "" || baseCurrency == "" {
			http.FailWithType(http.ParamError, c)
			return
		}
		dates := marketService.ListMarketStrikePriceOdAndDb(side, optionType, deliveryType, quoteCurrency, baseCurrency)
		http.OkWithData(dates, c)
		return
	}
}

func ListExpiryDatesByPrice(marketService core.MarketService, dbMarketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		price := c.Param("price")
		side := c.Query("side")
		optionType := c.Query("optionType")
		deliveryType := c.Query("deliveryType")
		quoteCurrency := c.Query("quoteCurrency")
		baseCurrency := c.Query("baseCurrency")
		if side == "" || optionType == "" || deliveryType == "" || quoteCurrency == "" || baseCurrency == "" {
			http.FailWithType(http.ParamError, c)
			return
		}
		dates, err := marketService.ListExpiryDatesByPriceOdAndDb(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(dates, c)
		return
	}
}

func ListDates(marketService core.MarketService, dbMarketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		side := c.Query("side")
		optionType := c.Query("optionType")
		deliveryType := c.Query("deliveryType")
		quoteCurrency := c.Query("quoteCurrency")
		baseCurrency := c.Query("baseCurrency")
		dates, err := marketService.ListDates(c, side, optionType, deliveryType, quoteCurrency, baseCurrency)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(dates, c)
		return
	}
}

func ListByDate(marketService core.MarketService, dbMarketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		dateMarkets, err := marketService.ListByDate(c, date)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(dateMarkets, c)
		return
	}
}

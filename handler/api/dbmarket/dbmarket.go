package dbmarket

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/pkg/http"
)

func FindByInstrument(marketService core.DbMarketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		instrument, err := marketService.ListMarketsByInstrument(name)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(instrument, c)
	}
}

func IndexPrice(propertyStore core.PropertyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		property, err := propertyStore.ReadProperty(c, core.DbBtcUsdIndexPrice)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(property, c)
	}
}

func ListStrikePrices(marketService core.DbMarketService) gin.HandlerFunc {
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
		dates := marketService.ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency)
		http.OkWithData(dates, c)
		return
	}
}

func ListExpiryDatesByPrice(marketService core.DbMarketService) gin.HandlerFunc {
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
		dates, err := marketService.ListExpiryDatesByPrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(dates, c)
		return
	}
}

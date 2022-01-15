package deliveryprice

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	http2 "option-dance/pkg/http"
)

func DeliveryPrice(store core.DeliveryPriceStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Query("date")
		asset := c.Query("asset")
		price, err := store.ReadPrice(c, asset, date)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		http2.OkWithData(price, c)
		return
	}
}

func DeliveryPrices(store core.DeliveryPriceStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		asset := c.Query("asset")
		prices, err := store.ListPricesByAsset(c, asset)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		http2.OkWithData(prices, c)
		return
	}
}

package statistics

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/pkg/http"
	"time"
)

func AnnualSellPutPremium(statisticsSrv core.StatisticsService, userSrv core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := userSrv.CurrentUserId(c)
		result, err := statisticsSrv.ListMonthlyPremium(c, uid, "ASK", time.Now().Year())
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(result, c)
		return
	}
}

func AnnualSellPutUnderlying(statisticsSrv core.StatisticsService, userSrv core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := userSrv.CurrentUserId(c)
		result, err := statisticsSrv.ListMonthlyUnderlying(c, uid, time.Now().Year())
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(result, c)
		return
	}
}

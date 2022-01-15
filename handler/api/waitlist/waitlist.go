package waitlist

import (
	"option-dance/core"
	"option-dance/pkg/http"

	"github.com/gin-gonic/gin"
)

func WaitListCreate(waitListSrv core.WaitListService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var waitlist core.WaitList
		err := c.BindJSON(&waitlist)
		if err != nil {
			http.FailWithMsg(err.Error(), c)
		}
		waitListSrv.Create(waitlist, c)
		http.Ok(c)
		return
	}
}

func WaitlistRankInfo(waitListSrv core.WaitListService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mid := c.Query("mid")
		info := waitListSrv.GetRankInfoByMid(mid)
		http.OkWithData(info, c)
		return
	}
}

func WaitlistWid(waitListSrv core.WaitListService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mid := c.Query("mid")
		email := c.Query("email")
		wid, _ := waitListSrv.GetWid(mid, email)
		http.OkWithData(wid, c)
		return
	}
}

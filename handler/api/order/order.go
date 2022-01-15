package order

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"option-dance/core"
	"option-dance/pkg/http"
	"option-dance/pkg/util"
	"option-dance/service/mtg"
	"time"
)

func Page(userz core.UserService, orderz core.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		order := c.Query("order")
		current, size, _, err := util.PageInfo(c)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		uid := userz.CurrentUserId(c)
		list, total, pages := orderz.Page(current, size, uid, status, order)
		http.OkWithData(gin.H{
			"records": list,
			"total":   total,
			"pages":   pages,
		}, c)
	}
}
func OrderCreateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var o core.OrderAction
		err := c.BindJSON(&o)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		//trace id check
		if len(o.TraceId) > 0 {
			_, err := uuid.FromString(o.TraceId)
			if err != nil {
				http.FailWithErr(err, c)
				return
			}
		} else {
			http.FailWithType(http.ParamError, c)
			return
		}
		payment, err := mtg.CreateOrderRequest(c, o)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(payment.Data, c)
		return
	}
}
func OrderCancel(userService core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("id")
		userId := userService.CurrentUserId(c)
		err := mtg.CancelOrderRequest(c, orderId, userId)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.Ok(c)
		return
	}
}
func OrderDetailByOrderId(userService core.UserService, orderz core.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("id")
		userId := userService.CurrentUserId(c)
		detail, err := orderz.FindByOrderId(c, userId, orderId)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(detail, c)
		return
	}
}

func OpenOrders(orderz core.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		detail, err := orderz.ListOpenOrdersBeforeExpiryDate(c, time.Now())
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(detail, c)
		return
	}
}

func OpenOrdersBeforeNow(orderz core.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		side := c.Query("side")
		if side == "" {
			side = core.PageSideAsk
		}
		detail, err := orderz.ListOpenOrdersBeforeTime(c, side, time.Now())
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(detail, c)
		return
	}
}

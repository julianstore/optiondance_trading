package position

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/pkg/http"
	"option-dance/pkg/util"
	"option-dance/service/mtg"
)

func ExerciseRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			InstrumentName string `json:"instrument_name"`
			Size           string `json:"size"`
			PositionId     string `json:"position_id"`
		}
		err := c.BindJSON(&req)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		result, err := mtg.ExerciseRequest(c, req.InstrumentName, req.Size, req.PositionId)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(result.Data, c)
		return
	}
}

func ExerciseCashDelivery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			InstrumentName string `json:"instrument_name"`
			Size           string `json:"size"`
			PositionId     string `json:"position_id"`
		}
		err := c.BindJSON(&req)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		_, err = mtg.CashDeliveryExerciseTx(c, req.InstrumentName, req.Size, req.PositionId)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.Ok(c)
		return
	}
}

func ListUserPositions(userService core.UserService, positionStore core.PositionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		order := c.Query("order")
		current, size, _, err := util.PageInfo(c)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		uid := userService.CurrentUserId(c)
		positions, total, pages, err := positionStore.ListUserPositions(current, size, uid, status, order)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(gin.H{
			"records": positions,
			"total":   total,
			"pages":   pages,
		}, c)
		return
	}
}

func PositionDetail(userService core.UserService, positionStore core.PositionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		uid := userService.CurrentUserId(c)
		position, err := positionStore.FindByPositionIdAndUid(c, uid, pid)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(position, c)
		return
	}
}

func DetailByID(userService core.UserService, positionService core.PositionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		uid := userService.CurrentUserId(c)
		position, err := positionService.PositionDetail(c, uid, pid)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(position, c)
		return
	}
}

func PositionDetailByInstrumentName(userService core.UserService, positionStore core.PositionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		uid := userService.CurrentUserId(c)
		position, err := positionStore.FindPositionByUidAndName(uid, name)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(position, c)
		return
	}
}

func PositionStatus(positionStore core.PositionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		status, err := positionStore.FindPositionStatusById(pid)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(status, c)
		return
	}
}

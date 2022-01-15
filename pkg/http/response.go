package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type Response struct {
	Code int32       `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS = 0
	FAILED  = 50000
)

func R(code int32, msg string) (r Response) {
	r = Response{
		Code: code,
		Msg:  msg,
	}
	return
}

func Result(code int32, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(FAILED, map[string]interface{}{}, "request failed", c)
}

func FailWithType(r Response, c *gin.Context) {
	Result(r.Code, map[string]interface{}{}, r.Msg, c)
	c.Abort()
}

func FailWithMsg(message string, c *gin.Context) {
	Result(FAILED, nil, message, c)
}

func FailWithErr(err error, c *gin.Context) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Result(404, nil, err.Error(), c)
		return
	}
	zap.L().Error("http response error:", zap.String("error", err.Error()))
	Result(FAILED, nil, err.Error(), c)
}

func FailWithDetailed(code int32, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}

func FailWithHttpStatus(httpStatus int, r Response, c *gin.Context) {
	c.JSON(httpStatus, Response{
		r.Code,
		r.Data,
		r.Msg,
	})
}

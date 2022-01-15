package hc

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Hc() gin.HandlerFunc {
	up := time.Now()
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
			"since":  time.Since(up).String(),
		})
	}
}

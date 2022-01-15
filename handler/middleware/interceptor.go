package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	response "option-dance/pkg/http"
	"time"
)

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 0 {
			payload := ParseJWT(token)
			if payload.ExpiredAt < time.Now().Unix() || len(payload.UserId) == 0 {
				c.Abort()
				response.FailWithHttpStatus(http.StatusUnauthorized, response.NotAuthorized, c)
			}
			c.Next()
		} else {
			c.Abort()
			response.FailWithHttpStatus(http.StatusUnauthorized, response.NotAuthorized, c)
		}
	}
}

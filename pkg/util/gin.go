package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func PageInfo(c *gin.Context) (current int64, size int64, qs string, err error) {
	currentStr := c.Query("current")
	sizeStr := c.Query("size")
	qs = c.Query("qs")
	cur, err := strconv.Atoi(currentStr)
	if err != nil {
		return current, size, qs, err
	}
	sz, err := strconv.Atoi(sizeStr)
	if err != nil {
		return current, size, qs, err
	}
	current = int64(cur)
	size = int64(sz)
	return
}

func Int64Param(c *gin.Context, key string) (p int64, err error) {
	v := c.Param(key)
	sz, err := strconv.Atoi(v)
	if err != nil {
		return p, err
	}
	p = int64(sz)
	return
}

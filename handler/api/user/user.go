package user

import (
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/pkg/http"
)

func HandlerUserMe(userSrv core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, err := userSrv.CurrentUser(c)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(currentUser, c)
		return
	}
}

func HandlerUserSettings(userSrv core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := userSrv.Settings(c)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.OkWithData(settings, c)
		return
	}
}

func HandlerUpdateUserSettings(userSrv core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto core.SettingsDTO
		err := c.BindJSON(&dto)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		err = userSrv.SettingsUpdate(dto, c)
		if err != nil {
			http.FailWithErr(err, c)
			return
		}
		http.Ok(c)
		return
	}
}

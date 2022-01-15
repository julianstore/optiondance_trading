package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"option-dance/core"
	"option-dance/handler/middleware"
	http2 "option-dance/pkg/http"
)

func HandleAuthByMixinToken(userService core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenDto struct {
			Token string `json:"token"`
		}
		err := c.BindJSON(&tokenDto)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		client := http.Client{}
		request, err := http.NewRequest("GET", "https://mixin-api.zeromesh.net/me", nil)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenDto.Token))
		resp, err := client.Do(request)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var userResp core.MixinUserResponse
		err = json.Unmarshal(body, &userResp)
		if err != nil {
			http2.FailWithErr(err, c)
			return
		}
		data := userResp.Data
		go func() {
			_, _ = userService.Save(&core.User{
				MixinUId: data.UserID,
				MixinID:  data.IdentityNumber,
				Nickname: data.FullName,
				Avatar:   data.AvatarURL,
				Type:     "MIXIN",
				Phone:    data.Phone,
			})
		}()
		jwt := middleware.GenerateJWT(data.UserID)
		http2.OkWithData(gin.H{
			"token": jwt,
		}, c)
	}
}

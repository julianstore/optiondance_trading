package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"option-dance/core"
	"option-dance/handler/middleware"
)

func New(userStore core.UserStore) core.UserService {
	return &service{userStore: userStore}
}

type service struct {
	userStore core.UserStore
}

func (s *service) Save(u *core.User) (user *core.User, err error) {
	user, err = s.userStore.Save(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) CurrentUserId(ctx *gin.Context) (userId string) {
	token := ctx.GetHeader("Authorization")
	payload := middleware.ParseJWT(token)
	userId = payload.UserId
	return
}

func (s *service) CurrentUser(c *gin.Context) (u *core.User, err error) {
	id := s.CurrentUserId(c)
	if len(id) == 0 {
		return nil, fmt.Errorf("dapp token maybe invalid")
	}
	u, err = s.userStore.FindByMixinId(id)
	if err != nil {
		return nil, err
	}
	return
}

func (s *service) Settings(c *gin.Context) (d *core.SettingsDTO, err error) {
	id := s.CurrentUserId(c)
	u, err := s.userStore.FindByMixinId(id)
	if err != nil {
		return nil, err
	}
	return &core.SettingsDTO{AppMode: u.AppMode, DeliveryType: u.DeliveryType}, nil
}

func (s *service) SettingsUpdate(dto core.SettingsDTO, c *gin.Context) (err error) {
	id := s.CurrentUserId(c)
	err = s.userStore.UpdateSettings(id, dto)
	if err != nil {
		return err
	}
	return nil
}

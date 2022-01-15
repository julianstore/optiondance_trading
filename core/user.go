package core

import (
	"github.com/gin-gonic/gin"
	"time"
)

type (
	User struct {
		ID           int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		MixinUId     string    `gorm:"column:mixin_uid;type:varchar(45)" json:"mixin_uid"`
		MixinID      string    `gorm:"column:mixin_id;type:varchar(45)" json:"mixin_id"`
		Nickname     string    `gorm:"column:nickname;type:varchar(45)" json:"nickname"`
		Avatar       string    `gorm:"column:avatar;type:varchar(200)" json:"avatar"`
		Type         string    `gorm:"column:type;type:varchar(45)" json:"type"`
		Email        string    `gorm:"column:email;type:varchar(45)" json:"email"`
		Phone        string    `gorm:"column:phone;type:varchar(45)" json:"phone"`
		Status       int8      `gorm:"column:status;type:tinyint" json:"status"`
		AppMode      int8      `gorm:"column:app_mode;type:tinyint" json:"app_mode"`
		DeliveryType int8      `gorm:"column:delivery_type;type:tinyint;default:0" json:"delivery_type"`
		CreatedAt    time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	}

	UserService interface {
		Save(u *User) (user *User, err error)
		CurrentUserId(ctx *gin.Context) (userId string)
		CurrentUser(c *gin.Context) (u *User, err error)
		Settings(c *gin.Context) (d *SettingsDTO, err error)
		SettingsUpdate(dto SettingsDTO, c *gin.Context) (err error)
	}

	UserStore interface {
		Save(u *User) (user *User, err error)
		FindByMixinId(mid string) (u *User, err error)
		UpdateSettings(mid string, dto SettingsDTO) error
	}
)

func (m *User) TableName() string {
	return "user"
}

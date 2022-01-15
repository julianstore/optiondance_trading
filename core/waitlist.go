package core

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	WaitList struct {
		ID          int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"` // Primary key
		Wid         int64     `gorm:"column:wid;type:bigint" json:"wid"`           // Invite wid
		MixinID     string    `gorm:"column:mixin_id;type:varchar(45)" json:"mixin_id"`
		MixinUId    string    `gorm:"column:mixin_uid;type:varchar(45)" json:"mixin_uid"`
		MixinName   string    `gorm:"column:mixin_name;type:varchar(45)" json:"mixin_name"`
		MixinPhone  string    `gorm:"column:mixin_phone;type:varchar(45)" json:"mixin_phone"`
		Email       string    `gorm:"column:email;type:varchar(45)" json:"email"`                    // Mail
		InviterWid  int64     `gorm:"column:inviter_wid;type:bigint" json:"inviter_wid"`             // Inviter wid
		InviteCount int64     `gorm:"column:invite_count;type:bigint;default:0" json:"invite_count"` // Number of invitations
		Type        int8      `gorm:"column:type;type:tinyint" json:"type"`                          // Type, 0: mixin, 1: other platforms
		CreatedAt   time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`          // Creation time
		UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`          // Update time
		Status      int8      `gorm:"column:status;type:tinyint" json:"status"`                      // state
		IP          string    `gorm:"column:ip;type:varchar(45)" json:"ip"`                          // User ip
		SystemInfo  string    `gorm:"column:system_info;type:varchar(45)" json:"system_info"`        // Device Information
	}

	WaitListRankInfo struct {
		Rank        int64 `json:"rank"`
		InviteCount int64 `json:"invite_count"`
	}

	WaitListService interface {
		Create(w WaitList, c *gin.Context)
		GetRankInfoByMid(mid string) (info WaitListRankInfo)
		GetWid(mid string, email string) (wid int64, err error)
		IncrInvitedCount(wid int64)
	}

	WaitListStore interface {
		Create(w WaitList) error
		GetByMixinId(mixinId string) (w *WaitList, err error)
		GetByWid(wid string) (w *WaitList, err error)
		GetByEmail(email string) (w *WaitList, err error)
		IncrInvitedCount(wid int64) error
	}
)

func (m *WaitList) TableName() string {
	return "waitlist"
}

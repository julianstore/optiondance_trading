package message

import (
	"context"
	"encoding/base64"
	"github.com/MixinNetwork/bot-api-go-client"
	"option-dance/cmd/config"
	"option-dance/core"
	"time"
)

func SaveMixinMsg(ctx context.Context, v MessageView) error {
	decodeData, err := base64.URLEncoding.DecodeString(v.Data)
	if err != nil {
		return err
	}

	message := core.MixinMessage{
		ConversationID:   v.ConversationId,
		UserID:           v.UserId,
		MessageID:        v.MessageId,
		Category:         v.Category,
		Data:             v.Data,
		DecodeData:       string(decodeData),
		RepresentativeID: v.RepresentativeId,
		QuoteMessageID:   v.QuoteMessageId,
		Status:           v.Status,
		Source:           v.Source,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
	err = config.Db.Model(core.MixinMessage{}).Create(&message).Error
	if err != nil {
		return err
	}
	//follow action
	m := config.Cfg.DApp
	if string(decodeData) == "你好" {
		user, err := bot.GetUser(ctx, v.UserId, m.AppID, m.SessionID, m.PrivateKey)
		if err != nil {
			return err
		}
		mixinUser := core.MixinUser{
			UserID:         user.UserId,
			Phone:          user.Phone,
			FullName:       user.FullName,
			AvatarURL:      user.AvatarURL,
			IdentityNumber: user.IdentityNumber,
			SessionID:      user.SessionId,
			PinToken:       user.PinToken,
			DeviceStatus:   user.DeviceStatus,
			CreatedAt:      time.Now(),
			UpdateTime:     time.Now(),
			Status:         0,
		}
		result := config.Db.Model(core.MixinUser{}).Where("user_id=?", user.UserId).Updates(&mixinUser)
		if result.RowsAffected == 0 {
			err = config.Db.Model(core.MixinUser{}).Create(&mixinUser).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

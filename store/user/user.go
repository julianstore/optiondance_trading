package user

import (
	"fmt"
	"gorm.io/gorm"
	"option-dance/core"
	"time"
)

func New(db *gorm.DB) core.UserStore {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (s *store) Save(u *core.User) (user *core.User, err error) {
	if len(u.MixinUId) == 0 {
		return nil, fmt.Errorf("SaveUserError: mixin uid should not be null")
	}
	var count int64
	if err = s.db.Model(&core.User{}).Where("mixin_uid = ?", u.MixinUId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		err := s.db.Model(&core.User{}).Create(&u).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.Model(&core.User{}).Where("mixin_uid = ?", u.MixinUId).Updates(&u).Error
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (s *store) FindByMixinId(mid string) (u *core.User, err error) {
	if err = s.db.Model(&core.User{}).Where("mixin_uid = ?", mid).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *store) UpdateSettings(mid string, dto core.SettingsDTO) error {
	if err := s.db.Model(&core.User{}).Where("mixin_uid = ?", mid).Updates(map[string]interface{}{
		"app_mode":      dto.AppMode,
		"delivery_type": dto.DeliveryType,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

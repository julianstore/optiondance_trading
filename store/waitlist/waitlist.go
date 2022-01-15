package waitlist

import (
	"gorm.io/gorm"
	"option-dance/core"
	"time"
)

func New(db *gorm.DB) core.WaitListStore {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (s *store) Create(w core.WaitList) error {
	if err := s.db.Model(core.WaitList{}).Create(&w).Error; err != nil {
		return nil
	}
	return nil
}

func (s *store) GetByMixinId(mixinId string) (w *core.WaitList, err error) {
	if err = s.db.Model(core.WaitList{}).Where("mixin_id=?", mixinId).First(&w).Error; err != nil {
		return nil, err
	}
	return w, nil
}

func (s *store) GetByWid(wid string) (w *core.WaitList, err error) {
	if err = s.db.Model(core.WaitList{}).Where("wid=?", wid).First(&w).Error; err != nil {
		return nil, err
	}
	return w, nil
}

func (s *store) GetByEmail(email string) (w *core.WaitList, err error) {
	if err = s.db.Model(core.WaitList{}).Where("email=?", email).First(&w).Error; err != nil {
		return nil, err
	}
	return w, nil
}

func (s *store) IncrInvitedCount(wid int64) (err error) {
	var w core.WaitList
	if err = s.db.Model(core.WaitList{}).Where("wid=?", wid).Find(&w).Error; err != nil {
		return err
	}
	if err = s.db.Model(core.WaitList{}).Where("id=?", wid).Updates(map[string]interface{}{
		"invite_count": w.InviteCount + 1,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

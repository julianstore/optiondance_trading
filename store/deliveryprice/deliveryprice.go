package deliveryprice

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	store2 "option-dance/store"
	"time"
)

func New(db *gorm.DB) core.DeliveryPriceStore {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (s *store) WritePrice(ctx context.Context, asset, date, price string, time time.Time) error {
	tx := s.db.Model(core.DeliveryPrice{}).Where("asset = ? AND date= ?", asset, date).Updates(map[string]interface{}{
		"price":      price,
		"updated_at": time,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		if err := s.db.Model(core.DeliveryPrice{}).Create(&core.DeliveryPrice{
			Date:      date,
			Price:     price,
			Asset:     asset,
			CreatedAt: time,
			UpdatedAt: time,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *store) ReadPrice(ctx context.Context, asset, date string) (string, error) {
	var deliveryPrice core.DeliveryPrice
	if err := s.db.Where("asset = ? and date = ?", asset, date).Find(&deliveryPrice).Error; store2.CheckErr(err) != nil {
		return "", err
	}
	return deliveryPrice.Price, nil
}

func (s *store) ListPricesByAsset(ctx context.Context, asset string) ([]*core.DeliveryPrice, error) {
	var deliveryPrices []*core.DeliveryPrice
	if err := s.db.Where("asset = ?", asset).Order("id desc").Limit(10).Find(&deliveryPrices).Error; store2.CheckErr(err) != nil {
		return nil, err
	}
	return deliveryPrices, nil
}

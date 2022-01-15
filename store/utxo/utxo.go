package utxo

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/store"
)

func New(db *gorm.DB) core.UtxoStore {
	return &utxoStore{db: db}
}

type utxoStore struct {
	db *gorm.DB
}

func (s *utxoStore) FindSpentBy(ctx context.Context, assetID, spentBy string) (utxo *core.UTXO, err error) {
	if err = s.db.Model(core.UTXO{}).WithContext(ctx).Where("asset_id = ? AND spent_by = ?", assetID, spentBy).
		Take(&utxo).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return utxo, nil
}

func (s *utxoStore) ListUnspentV2(ctx context.Context, assetId string, limit int) (utxos []*core.UTXO, err error) {

	if err := s.db.Model(core.UTXO{}).WithContext(ctx).
		Where("asset_id = ? AND spent_by = ?", assetId, "").
		Limit(limit).
		Order("id").
		Find(&utxos).Error; err != nil {
		return nil, err
	}
	for _, u := range utxos {
		_, err := core.WrapUTXORawData(u)
		if err != nil {
			return nil, nil
		}
	}
	return utxos, nil
}

func (s *utxoStore) Save(ctx context.Context, u *core.UTXO) (err error) {
	var oldUtxo core.UTXO
	if err := s.db.Model(core.UTXO{}).
		Where("trace_id = ?", u.TraceID).
		Select("id", "state").Find(&oldUtxo).Error; store.CheckErr(err) != nil {
		return err
	}
	if oldUtxo.ID == 0 {
		err := s.db.Model(core.UTXO{}).Create(&u).Error
		if err != nil {
			return err
		}
	} else {
		if oldUtxo.State != "spent" {
			if err := s.db.Model(core.UTXO{}).Where("id = ?", oldUtxo.ID).Updates(map[string]interface{}{
				"state":      u.State,
				"signed_by":  u.SignedBy,
				"signed_tx":  u.SignedTx,
				"json_data":  u.JSONData,
				"updated_at": u.UpdatedAt,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *utxoStore) List(ctx context.Context, lastId int64, limit int) (utxo []*core.UTXO, err error) {
	if err = s.db.Model(core.UTXO{}).WithContext(ctx).
		Where("id > ?", lastId).
		Limit(limit).
		Order("id").
		Find(&utxo).Error; err != nil {
		return nil, err
	}
	return utxo, nil
}

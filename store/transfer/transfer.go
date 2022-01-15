package transfer

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/store"
)

func New(db *gorm.DB) core.TransferStore {
	return &transferStore{db: db}
}

type transferStore struct {
	db *gorm.DB
}

func (s *transferStore) ListTransfers(ctx context.Context, status int8, limit int) (transfers []*core.Transfer, err error) {

	if err = s.db.WithContext(ctx).Where("handled = ?", status).
		Limit(limit).
		Order("id").
		Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

func (s *transferStore) UpdateTransfer(ctx context.Context, transfer *core.Transfer, status int8) error {
	if err := s.db.WithContext(ctx).Model(core.Transfer{}).Where("transfer_id = ?", transfer.TransferId).
		Update("handled", status).Error; err != nil {
		return err
	}
	return nil
}

func (s *transferStore) ReadTransferTrade(ctx context.Context, tradeId, assetId string) (trade *core.Trade, err error) {
	if err = s.db.WithContext(ctx).Model(core.Trade{}).Where("trade_id=?", tradeId).Find(&trade).Error; err != nil {
		return trade, err
	}
	if trade.FeeAssetID == assetId {
		return trade, nil
	}
	return nil, nil
}

func (s *transferStore) ListPendingTransfers(ctx context.Context, limit int) (transfers []*core.Transfer, err error) {

	if err = s.db.Model(core.Transfer{}).WithContext(ctx).
		Where("handled = 0").
		Limit(limit).
		Order("id").
		Find(&transfers).Error; err != nil {
		return nil, err
	}
	// filter by asset id
	filter := make(map[string]bool)
	var idx int

	for _, t := range transfers {
		if filter[t.AssetId] {
			continue
		}

		transfers[idx] = t
		filter[t.AssetId] = true
		idx++
	}

	transfers = transfers[:idx]
	return transfers, nil
}

func (s *transferStore) Create(ctx context.Context, transfer core.Transfer) error {
	if err := s.db.Model(core.Transfer{}).WithContext(ctx).Create(&transfer).Error; store.CheckErr(err) != nil {
		return err
	}
	return nil
}

func (s *transferStore) Exist(ctx context.Context, transferId string) (bool, error) {
	var count int64
	if err := s.db.Model(core.Transfer{}).WithContext(ctx).Where("transfer_id=?", transferId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

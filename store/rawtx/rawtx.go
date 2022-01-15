package rawtx

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/store"
)

func New(db *gorm.DB) core.RawTxStore {
	return &rawTxStore{db: db}
}

type rawTxStore struct {
	db *gorm.DB
}

func (s *rawTxStore) ListPendingRawTransactions(ctx context.Context, limit int) (result []*core.RawTransaction, err error) {
	if err = s.db.Model(core.RawTransaction{}).WithContext(ctx).Where("state = ?", "").Limit(limit).Find(&result).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return result, nil
}

func (s *rawTxStore) ExpireRawTransaction(ctx context.Context, tx *core.RawTransaction) error {
	return s.db.Model(core.RawTransaction{}).WithContext(ctx).Where("id = ?", tx.ID).Delete(tx).Error
}

func (s *rawTxStore) SaveRawTransaction(ctx context.Context, tx *core.RawTransaction) error {
	var count int64
	if err := s.db.Model(core.RawTransaction{}).WithContext(ctx).Where("trace_id = ?", tx.TraceID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return s.db.Model(core.RawTransaction{}).WithContext(ctx).Where("trace_id = ?", tx.TraceID).FirstOrCreate(tx).Error
	} else {
		return s.db.Model(core.RawTransaction{}).WithContext(ctx).Where("trace_id = ?", tx.TraceID).Updates(map[string]interface{}{
			"state":      tx.State,
			"error_info": tx.ErrorInfo,
		}).Error
	}
}

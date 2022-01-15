package transfer

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"gorm.io/gorm"
	"option-dance/cmd/config"
	"option-dance/core"
)

func New(transferStore core.TransferStore) core.TransferService {
	return &transferService{transferStore: transferStore}
}

type transferService struct {
	transferStore core.TransferStore
}

func (t *transferService) CreateRefundTransfer(ctx context.Context, s *core.UTXO) error {
	//amount := number.FromFloat(s.Amount).Mul(number.FromString("0.999"))
	amount := number.FromFloat(s.Amount)
	if amount.Exhausted() {
		return nil
	}

	transferId := core.GetSettlementId(s.TraceID, "REFUND")
	exist, err := t.transferStore.Exist(ctx, transferId)
	if err != nil {
		return err
	}

	if !exist {
		u, err := core.WrapUTXORawData(s)
		if err != nil {
			return err
		}
		var transferSource = core.TransferSourceUTXORefund
		action, err := core.MemoToAction(ctx, s)
		if err == nil {
			if action.AC == core.MsgPackActionCreateOrder {
				transferSource = core.TransferSourceOrderInvalid
			}
			if action.AC == core.MsgPackActionExercise {
				transferSource = core.TransferSourcePositionExercisedInvalid
			}
		}
		transferID := core.GetSettlementId(u.TraceID, "REFUND")
		var transfer = core.Transfer{
			TransferId: transferID,
			Source:     transferSource,
			Detail:     u.TraceID,
			AssetId:    u.AssetID,
			Amount:     amount.Persist(),
			CreatedAt:  s.CreatedAt,
			UserId:     u.SenderID,
			ClientId:   u.UserID,
			TraceID:    transferID,
			Opponents:  u.SenderID,
			Threshold:  1,
		}
		if err := t.transferStore.Create(ctx, transfer); err != nil {
			return err
		}
	}
	return nil
}

func (t *transferService) CreateClosePositionTransfer(tx *gorm.DB, margin float64, positionId, userId, optionType string, trade *core.Trade) error {

	amount := number.FromFloat(margin)
	if amount.Cmp(number.Zero()) == 0 {
		return nil
	}
	var (
		tradeId        = trade.TradeID
		transferSource = core.TransferSourcePositionClosedRefund
		detail         = positionId + ":" + tradeId
		transferId     = core.GetSettlementId(positionId, tradeId+":"+transferSource)
		dapp           = config.Cfg.DApp
		assetId        = trade.QuoteAssetID
	)
	if optionType == core.OptionTypeCALL {
		assetId = trade.BaseAssetID
	}
	var transfer = &core.Transfer{
		TransferId: transferId,
		Source:     transferSource,
		Detail:     detail,
		AssetId:    assetId,
		Amount:     amount.Persist(),
		CreatedAt:  trade.CreatedAt,
		UserId:     userId,
		ClientId:   dapp.AppID,
		TraceID:    transferId,
		Opponents:  userId,
		Threshold:  1,
	}
	if err := tx.Model(core.Transfer{}).Create(&transfer).Error; err != nil {
		return err
	}
	return nil
}

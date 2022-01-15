package position

import (
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"option-dance/cmd/config"
	"option-dance/core"
	"time"
)

// cash delivery transfers
func (p *positionService) CashBidPutPositionTransfers(exerciseBidPositions []*core.Position, updateDb bool,
	strikePrice, deliveryPrice number.Decimal, settledAt time.Time) (transferList []*core.Transfer, err error) {

	priceDiff := strikePrice.Sub(deliveryPrice)
	//priceDiff <= 0 indicate that exercise action is invalid，should be pass
	if priceDiff.LessThanOrEqual(decimal.Zero) {
		return transferList, nil
	}
	for _, e := range exerciseBidPositions {
		assetId := core.GetGlobalQuoteAssetId()
		amount := number.FromFloat(e.Size).Mul(priceDiff)
		//bid exercised transfer
		tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercised)
		transfer := core.Transfer{
			TransferId: tid,
			Source:     core.TransferSourcePositionExercised,
			Detail:     e.PositionID,
			AssetId:    assetId,
			Amount:     ToAmountString(decimal.RequireFromString(amount.String())),
			CreatedAt:  settledAt,
			UserId:     e.UserID,
			ClientId:   config.Cfg.DApp.AppID,
			TraceID:    tid,
			Opponents:  e.UserID,
			Threshold:  1,
		}
		transferList = append(transferList, &transfer)
		if updateDb {
			if err = p.positionStore.UpdateExercisedSize(e.ID, e.Size, settledAt); err != nil {
				return nil, err
			}
		}
	}
	return transferList, nil
}

func (p *positionService) CashAskPutPositionTransfers(askPositions []*core.Position, updateDb bool,
	strikePrice, deliveryPrice number.Decimal, exerciseRatio string, i core.OptionInfo, settledAt time.Time) (transferList []*core.Transfer, err error) {

	priceDiff := strikePrice.Sub(deliveryPrice)
	//priceDiff <= 0 indicate that exercise action is invalid，exerciseRatio should be set to zero
	if priceDiff.LessThanOrEqual(decimal.Zero) {
		exerciseRatio = "0"
	}

	for _, e := range askPositions {
		//use quote currency margin instead of underlying asset transfer
		underlyingAssetId, err := core.GetAssetIdByCurrency(i.QuoteCurrency)
		if err != nil {
			return nil, err
		}
		var exercisedSize string
		//no need to transfer underlying asset when exerciseRatio = 0
		if number.FromString(exerciseRatio).Cmp(number.Zero()) > 0 {
			amount := number.FromFloat(e.Size).Abs().
				Mul(deliveryPrice.Decimal).
				Mul(number.FromString(exerciseRatio).Abs()) //ask side use abs value
			exercisedSize = number.FromFloat(e.Size).Abs().String()
			tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercised)
			underlyingTransfer := core.Transfer{
				TransferId: tid,
				Source:     core.TransferSourcePositionExercised,
				Detail:     e.PositionID,
				AssetId:    underlyingAssetId,
				Amount:     ToAmountString(amount),
				CreatedAt:  settledAt,
				UserId:     e.UserID,
				ClientId:   config.Cfg.DApp.AppID,
				TraceID:    tid,
				Opponents:  e.UserID,
				Threshold:  1,
			}
			transferList = append(transferList, &underlyingTransfer)
		}

		// not exercised partial margin refund transfer
		quoteAssetId := core.GetGlobalQuoteAssetId()
		refundRatio := number.FromFloat(1).Sub(number.FromString(exerciseRatio))
		// no need to refund margin when refund ratio = 0
		if refundRatio.Cmp(number.FromString("0")) > 0 {
			remainingAmount := number.FromFloat(e.Size).Abs().
				Mul(number.FromString(refundRatio.String()).Abs()).
				Mul(strikePrice.Abs())
			tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercisedRefund)
			marginRefundTransfer := core.Transfer{
				TransferId: tid,
				Source:     core.TransferSourcePositionExercisedRefund,
				Detail:     e.PositionID,
				AssetId:    quoteAssetId,
				Amount:     ToAmountString(remainingAmount),
				CreatedAt:  settledAt,
				UserId:     e.UserID,
				ClientId:   config.Cfg.DApp.AppID,
				TraceID:    tid,
				Opponents:  e.UserID,
				Threshold:  1,
			}
			transferList = append(transferList, &marginRefundTransfer)
		}
		if updateDb {
			//update exercised size
			if len(exercisedSize) > 0 {
				size := number.FromString(exercisedSize).Float64()
				if err = p.positionStore.UpdateExercisedSize(e.ID, size, settledAt); err != nil {
					return nil, err
				}
			}
		}
		zap.L().Debug("askPosition settlement", zap.String("positionId", e.PositionID), zap.String("exercised_size", exercisedSize))
	}
	return transferList, nil
}

func (p *positionService) CashBidCallPositionTransfers(exerciseBidPositions []*core.Position, updateDb bool,
	strikePrice number.Decimal, optionInfo core.OptionInfo, settledAt time.Time) (transferList []*core.Transfer, err error) {
	//bid exercised user transfers
	for _, e := range exerciseBidPositions {
		assetId, _ := core.GetAssetIdByCurrency(optionInfo.BaseCurrency)
		amount := number.FromFloat(e.Size)
		//bid exercised transfer
		tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercised)
		transfer := core.Transfer{
			TransferId: tid,
			Source:     core.TransferSourcePositionExercised,
			Detail:     e.PositionID,
			AssetId:    assetId,
			Amount:     ToAmountString(decimal.RequireFromString(amount.String())),
			CreatedAt:  settledAt,
			UserId:     e.UserID,
			ClientId:   config.Cfg.DApp.AppID,
			TraceID:    tid,
			Opponents:  e.UserID,
			Threshold:  1,
		}
		transferList = append(transferList, &transfer)
		if updateDb {
			if err = p.positionStore.UpdateExercisedSize(e.ID, e.Size, settledAt); err != nil {
				return nil, err
			}
		}
	}
	return transferList, nil
}

func (p *positionService) CashAskCallPositionTransfers(askPositions []*core.Position, updateDb bool,
	strikePrice number.Decimal, exerciseRatio string, i core.OptionInfo, settledAt time.Time) (transferList []*core.Transfer, err error) {
	//ask exercised user transfers
	for _, e := range askPositions {
		//quote asset transfer
		quoteAssetId := core.GetGlobalQuoteAssetId()
		var exercisedSize string
		//no need to transfer quote asset when exerciseRatio = 0
		if number.FromString(exerciseRatio).Cmp(number.Zero()) > 0 {
			amount := number.FromFloat(e.Size).Abs().Mul(strikePrice.Decimal).
				Mul(number.FromString(exerciseRatio).Abs()) //ask side use abs value
			exercisedSize = amount.String()
			tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercised)
			quoteAssetTransfer := core.Transfer{
				TransferId: tid,
				Source:     core.TransferSourcePositionExercised,
				Detail:     e.PositionID,
				AssetId:    quoteAssetId,
				Amount:     ToAmountString(amount),
				CreatedAt:  settledAt,
				UserId:     e.UserID,
				ClientId:   config.Cfg.DApp.AppID,
				TraceID:    tid,
				Opponents:  e.UserID,
				Threshold:  1,
			}
			transferList = append(transferList, &quoteAssetTransfer)
		}

		// not exercised partial margin underlying refund transfer
		underlyingAssetId, err := core.GetAssetIdByCurrency(i.BaseCurrency)
		if err != nil {
			return nil, err
		}
		refundRatio := number.FromFloat(1).Sub(number.FromString(exerciseRatio))
		// no need to refund margin when refund ratio = 0
		if refundRatio.Cmp(number.Zero()) > 0 {
			remainingAmount := number.FromFloat(e.Size).Abs().
				Mul(number.FromString(refundRatio.String()).Abs())
			tid := core.GetSettlementId(e.PositionID, core.TransferSourcePositionExercisedRefund)
			marginRefundTransfer := core.Transfer{
				TransferId: tid,
				Source:     core.TransferSourcePositionExercisedRefund,
				Detail:     e.PositionID,
				AssetId:    underlyingAssetId,
				Amount:     ToAmountString(remainingAmount),
				CreatedAt:  settledAt,
				UserId:     e.UserID,
				ClientId:   config.Cfg.DApp.AppID,
				TraceID:    tid,
				Opponents:  e.UserID,
				Threshold:  1,
			}
			transferList = append(transferList, &marginRefundTransfer)
		}
		if updateDb {
			//update exercised size
			if len(exercisedSize) > 0 {
				size := number.FromString(exercisedSize).Float64()
				if err = p.positionStore.UpdateExercisedSize(e.ID, size, settledAt); err != nil {
					return nil, err
				}
			}
		}
		zap.L().Debug("askPosition settlement", zap.String("positionId", e.PositionID), zap.String("exercised_size", exercisedSize))
	}
	return transferList, nil
}

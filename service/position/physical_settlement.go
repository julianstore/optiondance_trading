package position

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"option-dance/cmd/config"
	"option-dance/core"
	"time"
)

// calculate bid positions exercise ratio
func (p *positionService) CalcExerciseRatio(ctx context.Context, instrumentName, deliveryType string,
	deliveryPrice, strikePrice number.Decimal, settledAt time.Time, autoExercise bool) (exerciseBidPositions []*core.Position, ratio string, err error) {

	var (
		bidTotalSize          decimal.Decimal
		exerciseRequestedSize decimal.Decimal
		allBidPositions       []*core.Position
	)

	ratio = "0"

	//total bid positions
	AllStatusSet := []int64{core.PositionStatusNotExercised, core.PositionStatusExerciseRequested}
	allBidPositions, err = p.positionStore.ListSidePositions(ctx, instrumentName, "BID", AllStatusSet)
	if err != nil {
		return nil, ratio, err
	}
	for _, e := range allBidPositions {
		bidTotalSize = bidTotalSize.Add(number.FromFloat(e.Size).Decimal)
	}

	// Exercise Requested bid positions
	if deliveryType == core.DeliveryTypeCash && autoExercise {
		if deliveryPrice.LessThan(strikePrice.Decimal) {
			ratio = "1"
			exerciseBidPositions = allBidPositions
			if err = p.positionStore.AutoExercisePositionWithName(ctx, instrumentName, settledAt); err != nil {
				return nil, ratio, err
			}
		}
	} else {
		statusSet := []int64{core.PositionStatusExerciseRequested}
		exerciseBidPositions, err = p.positionStore.ListSidePositions(ctx, instrumentName, "BID", statusSet)
		if err != nil {
			return nil, ratio, err
		}
		for _, e := range exerciseBidPositions {
			exerciseRequestedSize = exerciseRequestedSize.Add(number.FromFloat(e.Size).Decimal)
		}

		//exerciseRatio
		if bidTotalSize.Cmp(decimal.Zero) != 0 {
			ratio = exerciseRequestedSize.Div(bidTotalSize).String()
		}
	}
	return
}

// physical delivery transfers
func (p *positionService) PhysicalBidPutPositionTransfers(exerciseBidPositions []*core.Position, updateDb bool,
	strikePrice number.Decimal, settledAt time.Time) (transferList []*core.Transfer, err error) {
	//bid exercised user transfers
	for _, e := range exerciseBidPositions {
		assetId := core.GetGlobalQuoteAssetId()
		amount := number.FromFloat(e.Size).Mul(strikePrice)
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

func (p *positionService) PhysicalAskPutPositionTransfers(askPositions []*core.Position, updateDb bool,
	strikePrice number.Decimal, exerciseRatio string, i core.OptionInfo, settledAt time.Time) (transferList []*core.Transfer, err error) {
	//ask exercised user transfers
	for _, e := range askPositions {
		//underlying asset transfer
		underlyingAssetId, err := core.GetAssetIdByCurrency(i.BaseCurrency)
		if err != nil {
			return nil, err
		}
		var exercisedSize string
		//no need to transfer underlying asset when exerciseRatio = 0
		if number.FromString(exerciseRatio).Cmp(number.FromString("0")) > 0 {
			amount := number.FromFloat(e.Size).Abs().Mul(number.FromString(exerciseRatio).Abs()) //ask side use abs value
			exercisedSize = amount.String()
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

func (p *positionService) PhysicalBidCallPositionTransfers(exerciseBidPositions []*core.Position, updateDb bool,
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

func (p *positionService) PhysicalAskCallPositionTransfers(askPositions []*core.Position, updateDb bool,
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

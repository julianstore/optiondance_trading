package trade

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/match/engine"
	"time"
)

const (
	MakerFeeRate = "0.000"
	TakerFeeRate = "0.000"

	LiquidityTaker = "TAKER"
	LiquidityMaker = "MAKER"
)

func NewTradeService(notifier core.Notifier,
	positionService core.PositionService, orderStore core.OrderStore) core.TradeService {
	return &tradeService{notifier: notifier, positionService: positionService, orderStore: orderStore}
}

type tradeService struct {
	orderStore      core.OrderStore
	notifier        core.Notifier
	positionService core.PositionService
}

func (s *tradeService) Transact(ctx context.Context, taker, maker *core.EngineOrder, amount number.Integer) (string, error) {
	askTrade, bidTrade := s.MakeTrades(taker, maker, amount.Decimal())
	askTransfer, bidTransfer, bidRefundTransfer := s.HandleFeesWithBidRefund(askTrade, bidTrade, taker, maker)
	if err := config.Db.Transaction(func(tx *gorm.DB) error {
		var askTransferCount int64
		if err := tx.Model(core.Transfer{}).Where("trace_id=?", askTransfer.TraceID).Count(&askTransferCount).Error; err != nil {
			return err
		}
		if askTransferCount == 0 {
			if err := tx.Model(core.Transfer{}).Create(&askTransfer).Error; err != nil {
				return err
			}
		}
		if bidRefundTransfer != nil && number.FromString(bidRefundTransfer.Amount).Cmp(number.Zero()) > 0 {
			var bidRefundTransferCount int64
			if err := tx.Model(core.Transfer{}).Where("trace_id=?", bidRefundTransfer.TraceID).Count(&bidRefundTransferCount).Error; err != nil {
				return err
			}
			if bidRefundTransferCount == 0 {
				if err := tx.Model(core.Transfer{}).Create(&bidRefundTransfer).Error; err != nil {
					return err
				}
			}
		}
		go func() {
			if err := s.notifier.Trade(ctx, bidTrade, bidTransfer); err != nil {
				zap.L().Error("messageStore.create", zap.String("type", "GetBizInfoMsgByTrade"))
			}
		}()

		if err := tx.Model(core.Trade{}).Create(&askTrade).Error; err != nil {
			return err
		}
		if err := tx.Model(core.Trade{}).Create(&bidTrade).Error; err != nil {
			return err
		}
		if err := s.positionService.UpdatePositionWithTrade(tx, askTrade); err != nil {
			return err
		}
		if err := s.positionService.UpdatePositionWithTrade(tx, bidTrade); err != nil {
			return err
		}
		if err := s.MakeOrderMutations(taker, maker, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
	return askTrade.TradeID, nil
}

// CancelOrder cancel order
func (s *tradeService) CancelOrder(ctx context.Context, order *core.EngineOrder, cancelledAt time.Time) error {

	oldOrder, err := s.orderStore.FindByOrderId(ctx, order.Id)
	if err != nil {
		return err
	}
	cancelOrderTransfer := true
	if oldOrder == nil || oldOrder.OrderStatus > core.OrderStatusCanceling {
		cancelOrderTransfer = false
	}

	orderUpdateStruct := core.Order{
		OrderID:         order.Id,
		RemainingAmount: order.RemainingAmount.Persist(),
		FilledAmount:    order.FilledAmount.Persist(),
		RemainingFunds:  order.RemainingFunds.Persist(),
		FilledFunds:     order.FilledFunds.Persist(),
		OrderStatus:     core.OrderStatusCancelled,
	}
	err = config.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(core.Order{}).
			WithContext(ctx).Where("order_id=?", order.Id).Updates(orderUpdateStruct).Error
		if err != nil {
			return err
		}
		err = tx.Model(core.Action{}).WithContext(ctx).
			Where("order_id=? and action in (?)", order.Id, []string{engine.OrderActionCreate, engine.OrderActionCancel}).
			Delete(core.Order{}).Error
		if err != nil {
			return err
		}
		//handle cancel order transfer
		transfer := &core.Transfer{
			TransferId: core.GetSettlementId(order.Id, engine.OrderActionCancel),
			Source:     core.TransferSourceOrderCancelled,
			Detail:     order.Id,
			AssetId:    order.Base,
			Amount:     order.RemainingAmount.Persist(),
			CreatedAt:  cancelledAt,
			UserId:     order.UserId,
			ClientId:   order.BrokerId,
			TraceID:    order.Id,
			Opponents:  order.UserId,
			Threshold:  1,
		}
		// sell put
		remainingAmount := order.RemainingAmount.Decimal()
		filledAmount := order.FilledAmount.Decimal()
		totalAmount := remainingAmount.Add(filledAmount)
		remainingFunds := order.RemainingFunds.Decimal()
		filledFunds := order.FilledFunds.Decimal()
		totalFunds := remainingFunds.Add(filledFunds)
		var remainingRatio number.Decimal
		if order.Side == "ASK" {
			remainingRatio = remainingAmount.Div(totalAmount)
		}
		if order.Side == "BID" {
			remainingRatio = remainingFunds.Div(totalFunds)
		}
		//sell put
		if order.OptionType == core.OptionTypePUT && order.Side == "ASK" {
			transfer.AssetId = order.Quote
			transfer.Amount = number.FromString(order.Margin).Mul(remainingRatio).Round(8).Persist()
		}
		//buy put
		if order.OptionType == core.OptionTypePUT && order.Side == "BID" {
			transfer.AssetId = order.Quote
			transfer.Amount = order.RemainingFunds.Persist()
		}
		//sell call
		if order.OptionType == core.OptionTypeCALL && order.Side == "ASK" {
			transfer.AssetId = order.Base
			transfer.Amount = number.FromString(order.Margin).Mul(remainingRatio).Round(8).Persist()
		}
		//buy call
		if order.OptionType == core.OptionTypeCALL && order.Side == "BID" {
			transfer.AssetId = order.Quote
			transfer.Amount = order.RemainingFunds.Persist()
		}
		if cancelOrderTransfer {
			err = tx.Model(&core.Transfer{}).WithContext(ctx).Create(transfer).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *tradeService) MakeOrderMutations(taker, maker *core.EngineOrder, tx *gorm.DB) (err error) {
	takerUpdateStruct := core.Order{
		OrderID:         taker.Id,
		RemainingAmount: taker.RemainingAmount.Persist(),
		FilledAmount:    taker.FilledAmount.Persist(),
		RemainingFunds:  taker.RemainingFunds.Persist(),
		FilledFunds:     taker.FilledFunds.Persist(),
	}
	makerUpdateStruct := core.Order{
		OrderID:         maker.Id,
		RemainingAmount: maker.RemainingAmount.Persist(),
		FilledAmount:    maker.FilledAmount.Persist(),
		RemainingFunds:  maker.RemainingFunds.Persist(),
		FilledFunds:     maker.FilledFunds.Persist(),
	}
	if taker.Filled() {
		takerUpdateStruct.OrderStatus = core.OrderStateDone
	}
	if maker.Filled() {
		makerUpdateStruct.OrderStatus = core.OrderStateDone
	}

	err = tx.Model(&core.Order{}).Where("order_id=?", maker.Id).Updates(makerUpdateStruct).Error
	if err != nil {
		return err
	}
	err = tx.Model(&core.Order{}).Where("order_id=?", taker.Id).Updates(takerUpdateStruct).Error
	if err != nil {
		return err
	}

	if taker.Filled() {
		err = tx.Model(core.Action{}).Where("order_id=? and action = ?", taker.Id, engine.OrderActionCreate).Delete(core.Action{}).Error
		if err != nil {
			return err
		}
		err = tx.Model(core.Action{}).Where("order_id=? and action = ?", taker.Id, engine.OrderActionCancel).Delete(core.Action{}).Error
		if err != nil {
			return err
		}
	}
	if maker.Filled() {
		err = tx.Model(core.Action{}).Where("order_id=? and action = ?", maker.Id, engine.OrderActionCreate).Delete(core.Action{}).Error
		if err != nil {
			return err
		}
		err = tx.Model(core.Action{}).Where("order_id=? and action = ?", maker.Id, engine.OrderActionCancel).Delete(core.Action{}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *tradeService) MakeTrades(taker, maker *core.EngineOrder, amount number.Decimal) (*core.Trade, *core.Trade) {
	tradeId := core.GetSettlementId(taker.Id, maker.Id)
	askOrderId, bidOrderId := taker.Id, maker.Id
	if taker.Side == engine.PageSideBid {
		askOrderId, bidOrderId = maker.Id, taker.Id
	}
	price := maker.Price.Decimal()

	takerTrade := &core.Trade{
		TradeID:        tradeId,
		Liquidity:      LiquidityTaker,
		AskOrderID:     askOrderId,
		BidOrderID:     bidOrderId,
		QuoteAssetID:   taker.Quote,
		BaseAssetID:    taker.Base,
		Side:           taker.Side,
		Price:          price.Persist(),
		Amount:         amount.Persist(),
		CreatedAt:      taker.CreatedAt,
		UserID:         taker.UserId,
		InstrumentName: taker.InstrumentName,
		Status:         core.TradeStatusOpen,
	}
	makerTrade := &core.Trade{
		TradeID:        tradeId,
		Liquidity:      LiquidityMaker,
		AskOrderID:     askOrderId,
		BidOrderID:     bidOrderId,
		QuoteAssetID:   maker.Quote,
		BaseAssetID:    maker.Base,
		Side:           maker.Side,
		Price:          price.Persist(),
		Amount:         amount.Persist(),
		CreatedAt:      taker.CreatedAt,
		UserID:         maker.UserId,
		InstrumentName: maker.InstrumentName,
		Status:         core.TradeStatusOpen,
	}

	askTrade, bidTrade := takerTrade, makerTrade
	if askTrade.Side == engine.PageSideBid {
		askTrade, bidTrade = makerTrade, takerTrade
	}
	return askTrade, bidTrade
}

func (s *tradeService) HandleFees(ask, bid *core.Trade, taker, maker *core.EngineOrder) (*core.Transfer, *core.Transfer) {
	total := number.FromString(ask.Amount).Mul(number.FromString(ask.Price))
	askFee := total.Mul(number.FromString(TakerFeeRate))
	bidFee := number.FromString(bid.Amount).Mul(number.FromString(MakerFeeRate))
	if ask.Liquidity == LiquidityMaker {
		askFee = total.Mul(number.FromString(MakerFeeRate))
		bidFee = number.FromString(bid.Amount).Mul(number.FromString(TakerFeeRate))
	}

	ask.FeeAssetID = ask.QuoteAssetID
	ask.FeeAmount = askFee.Persist()
	bid.FeeAssetID = bid.BaseAssetID
	bid.FeeAmount = bidFee.Persist()

	askTransfer := &core.Transfer{
		TransferId: core.GetSettlementId(ask.TradeID, ask.Liquidity),
		Source:     core.TransferSourceTradeConfirmed,
		Detail:     ask.TradeID,
		AssetId:    ask.QuoteAssetID,
		Amount:     total.Sub(askFee).Persist(),
		CreatedAt:  taker.CreatedAt,
		UserId:     ask.UserID,
		TraceID:    ask.TradeID,
		Threshold:  1,
		Opponents:  ask.UserID,
	}
	bidTransfer := &core.Transfer{
		TransferId: core.GetSettlementId(bid.TradeID, bid.Liquidity),
		Source:     core.TransferSourceTradeConfirmed,
		Detail:     bid.TradeID,
		AssetId:    bid.QuoteAssetID,
		Amount:     number.FromString(bid.Amount).Sub(bidFee).Persist(),
		CreatedAt:  taker.CreatedAt,
		UserId:     bid.UserID,
		TraceID:    bid.TradeID,
		Threshold:  1,
		Opponents:  bid.UserID,
	}
	if taker.Side == engine.PageSideAsk {
		askTransfer.ClientId = maker.BrokerId
		bidTransfer.ClientId = taker.BrokerId
	} else {
		askTransfer.ClientId = taker.BrokerId
		bidTransfer.ClientId = maker.BrokerId
	}
	return askTransfer, bidTransfer
}

func (s *tradeService) HandleFeesWithBidRefund(ask, bid *core.Trade, taker, maker *core.EngineOrder) (*core.Transfer, *core.Transfer, *core.Transfer) {
	total := number.FromString(ask.Amount).Mul(number.FromString(ask.Price))
	askFee := total.Mul(number.FromString(TakerFeeRate))
	bidFee := number.FromString(bid.Amount).Mul(number.FromString(MakerFeeRate))
	if ask.Liquidity == LiquidityMaker {
		askFee = total.Mul(number.FromString(MakerFeeRate))
		bidFee = number.FromString(bid.Amount).Mul(number.FromString(TakerFeeRate))
	}

	ask.FeeAssetID = ask.QuoteAssetID
	ask.FeeAmount = askFee.Persist()
	bid.FeeAssetID = bid.BaseAssetID
	bid.FeeAmount = bidFee.Persist()

	askTransferID := core.GetSettlementId(ask.TradeID, ask.Liquidity)
	askTransfer := &core.Transfer{
		TransferId: askTransferID,
		Source:     core.TransferSourceTradeConfirmed,
		Detail:     ask.TradeID,
		AssetId:    ask.QuoteAssetID,
		Amount:     total.Sub(askFee).Persist(),
		CreatedAt:  taker.CreatedAt,
		UserId:     ask.UserID,
		TraceID:    askTransferID,
		Threshold:  1,
		Opponents:  ask.UserID,
	}
	bidTransferID := core.GetSettlementId(bid.TradeID, bid.Liquidity)
	bidTransfer := &core.Transfer{
		TransferId: bidTransferID,
		Source:     core.TransferSourceTradeConfirmed,
		Detail:     bid.TradeID,
		AssetId:    bid.QuoteAssetID,
		Amount:     number.FromString(bid.Amount).Sub(bidFee).Persist(),
		CreatedAt:  taker.CreatedAt,
		UserId:     bid.UserID,
		TraceID:    bidTransferID,
		Threshold:  1,
		Opponents:  bid.UserID,
	}
	if taker.Side == engine.PageSideAsk {
		askTransfer.ClientId = maker.BrokerId
		bidTransfer.ClientId = taker.BrokerId
	} else {
		askTransfer.ClientId = taker.BrokerId
		bidTransfer.ClientId = maker.BrokerId
	}

	//handle bid refund
	var bidRefundTransfer *core.Transfer
	var refundBidOrder *core.EngineOrder

	if taker.Side == engine.PageSideBid && taker.Filled() {
		refundBidOrder = taker
	}
	if maker.Side == engine.PageSideBid && maker.Filled() {
		refundBidOrder = maker
	}
	if refundBidOrder != nil {
		bidRefundTransferID := core.GetSettlementId(refundBidOrder.Id, core.TransferSourceTradeBidRefund)
		amount := refundBidOrder.RemainingFunds
		bidRefundTransfer = &core.Transfer{
			TransferId: bidRefundTransferID,
			Source:     core.TransferSourceTradeBidRefund,
			Detail:     refundBidOrder.Id,
			AssetId:    refundBidOrder.Quote,
			Amount:     amount.Persist(),
			CreatedAt:  refundBidOrder.CreatedAt,
			UserId:     refundBidOrder.UserId,
			TraceID:    bidRefundTransferID,
			Threshold:  1,
			Opponents:  refundBidOrder.UserId,
			ClientId:   refundBidOrder.BrokerId,
		}
	}
	return askTransfer, bidTransfer, bidRefundTransfer
}

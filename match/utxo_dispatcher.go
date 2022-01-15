package match

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"option-dance/core"
	"option-dance/match/engine"
	time2 "option-dance/pkg/time"
	"option-dance/service/order"
	"strconv"
	"time"
)

const (
	CheckpointUtxo = "utxo_process_checkpoint"
)

var (
	gasAmount = decimal.NewFromFloat(0.00000001)
)

type UtxoDispatcher struct {
	propertyStore      core.PropertyStore
	positionStore      core.PositionStore
	positionService    core.PositionService
	orderStore         core.OrderStore
	utxoStore          core.UtxoStore
	orderService       core.OrderService
	transferService    core.TransferService
	deliveryPriceStore core.DeliveryPriceStore
	notifier           core.Notifier
	actionMap          map[string]func(ctx context.Context, u *core.UTXO, action *core.MsgPack) error
}

func NewUtxoDispatcher(
	propertyStore core.PropertyStore,
	positionStore core.PositionStore,
	orderStore core.OrderStore,
	orderService core.OrderService,
	utxoStore core.UtxoStore,
	positionService core.PositionService,
	transferService core.TransferService,
	notifier core.Notifier,
	deliveryPriceStore core.DeliveryPriceStore,
) UtxoDispatcher {
	d := UtxoDispatcher{
		propertyStore:      propertyStore,
		positionStore:      positionStore,
		notifier:           notifier,
		orderStore:         orderStore,
		orderService:       orderService,
		positionService:    positionService,
		utxoStore:          utxoStore,
		transferService:    transferService,
		deliveryPriceStore: deliveryPriceStore,
	}
	d.actionMap = map[string]func(ctx context.Context, u *core.UTXO, action *core.MsgPack) error{
		core.MsgPackActionCreateOrder:         d.createOrderHandler,
		core.MsgPackActionCancelOrder:         d.cancelOrderHandler,
		core.MsgPackActionExercise:            d.exerciseHandler,
		core.MsgPackActionAutoExercise:        d.autoExerciseHandler,
		core.MsgPackActionExpiringCancelOrder: d.ExpiringCancelOrder,
		core.MsgPackActionExpiringNotify:      d.ExpiryNotify,
		core.MsgPackActionSettlement:          d.Settlement,
		core.MsgPackSyncDbDeliveryPrice:       d.SyncDbDeliveryPrice,
	}
	return d
}

func (d *UtxoDispatcher) Run(ctx context.Context) error {
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := d.process(ctx); err == nil {
				dur = 150 * time.Millisecond
			} else {
				zap.L().Error("LoopingDispatchUTXO error", zap.Error(err))
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (d *UtxoDispatcher) process(ctx context.Context) error {
	const limit = 500
	checkpoint, err := d.propertyStore.ReadProperty(ctx, CheckpointUtxo)
	if err != nil {
		return err
	}
	var lastId = 0
	if checkpoint != "" {
		lastId, err = strconv.Atoi(checkpoint)
		if err != nil {
			return err
		}
	}
	utxos, err := d.utxoStore.List(ctx, int64(lastId), limit)
	if err != nil {
		return err
	}
	for _, s := range utxos {
		d.ensureProcessUTXO(ctx, s)
		checkpoint = strconv.Itoa(int(s.ID))
	}
	err = d.propertyStore.WriteProperty(ctx, CheckpointUtxo, checkpoint)
	if err != nil {
		return err
	}
	return nil
}

func (d *UtxoDispatcher) ensureProcessUTXO(ctx context.Context, s *core.UTXO) {
	for {
		err := d.processUTXO(ctx, s)
		if err == nil {
			break
		}
		zap.L().Error("ensureProcessSnapshot", zap.Error(err))
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *UtxoDispatcher) processUTXO(ctx context.Context, u *core.UTXO) error {
	if u.SenderID == "" || u.TraceID == "" {
		return nil
	}
	if u.Amount <= 0 {
		return nil
	}
	action, err := core.MemoToAction(ctx, u)
	if err != nil {
		zap.L().Error("parse memo to action error", zap.Error(err))
		return d.refundUTXO(ctx, u)
	}
	if fn, ok := d.actionMap[action.AC]; ok {
		if err = fn(ctx, u, action); err != nil {
			return err
		}
	}
	return nil
}

// create order handler
func (d *UtxoDispatcher) createOrderHandler(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {

	zap.L().Info("Create Order", zap.String("orderId", action.TC))

	err := action.CheckCreateOrderPackWithUtxo(u)
	if err != nil {
		zap.L().Error("checkOrderMsgPack error", zap.Error(err))
		return d.refundUTXO(ctx, u)
	}

	quote, base := core.GetQuoteBasePair(action.QC, action.BC)
	if quote == "" || base == "" {
		zap.L().Error("quote base currency not supported: " + quote + ":" + base)
		return d.refundUTXO(ctx, u)
	}

	price := number.FromString(action.P)
	optionSize := number.FromString(action.A)
	remainingFunds := number.Zero()
	remainingAmount := optionSize
	if action.S == engine.PageSideBid {
		remainingFunds = optionSize.Mul(price)
	} else {
		remainingAmount = optionSize
		if action.T == core.OrderTypeLimit && price.Mul(remainingAmount).Cmp(QuoteMinimum(quote)) < 0 {
			return d.refundUTXO(ctx, u)
		}
	}
	err = d.orderService.CreateOrderAction(ctx, &core.Order{
		OrderID:         action.TC,
		Side:            action.S,
		OrderType:       action.T,
		Price:           action.P,
		RemainingAmount: remainingAmount.Persist(),
		FilledAmount:    number.Zero().Persist(),
		RemainingFunds:  remainingFunds.Persist(),
		FilledFunds:     number.Zero().Persist(),
		Margin:          action.M,
		QuoteAssetID:    quote,
		BaseAssetID:     base,
		InstrumentName:  action.I,
		QuoteCurrency:   action.QC,
		BaseCurrency:    action.BC,
	}, u.SenderID, order.OrderActionTypeUser, u.CreatedAt)
	return err
}

// exercise Handler
func (d *UtxoDispatcher) exerciseHandler(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {

	size := action.A
	positionId := action.TC
	instrumentName := action.I

	zap.L().Info("position exercised request", zap.String("positionId", positionId))

	//time check
	instrument, err := core.ParseInstrument(instrumentName)
	if err != nil {
		zap.L().Info("exercise:name format error", zap.String("positionId", positionId), zap.String("name", instrumentName))
		return d.refundUTXO(ctx, u)
	}

	optionType := instrument.OptionType
	expiryTs := instrument.ExpirationTimestamp
	exercisedTs := u.CreatedAt.Unix()

	if !core.ExercisableContract(exercisedTs, expiryTs) {
		zap.L().Warn("exercised at a wrong time", zap.String("positionId", positionId), zap.Time("exercisedAt", u.CreatedAt))
		return d.refundUTXO(ctx, u)
	}

	position, err := d.positionStore.FindByPositionId(ctx, positionId)
	if err != nil {
		return err
	}
	if position.ID == 0 || position.Status != core.PositionStatusNotExercised {
		zap.L().Info("exercise:wrong position status", zap.String("positionId", positionId), zap.Int8("status", position.Status))
		return d.refundUTXO(ctx, u)
	}

	if instrumentName != position.InstrumentName {
		zap.L().Info("exercise:wrong position name", zap.String("positionId", positionId), zap.String("name", position.InstrumentName))
		return d.refundUTXO(ctx, u)
	}

	// cash delivery exercise
	if instrument.DeliveryType == core.DeliveryTypeCash {
		if err := d.positionStore.ExercisePositionWithID(ctx, positionId, instrumentName, u.CreatedAt); err != nil {
			zap.L().Error("positionStore.positionExercised", zap.Error(err))
			return err
		}
		return nil
	}

	//physical delivery position exercise check
	baseAsset, err := core.GetAssetIdByCurrency(position.BaseCurrency)
	if err != nil {
		zap.L().Error("exercise:GetAssetIdByCurrency", zap.String("positionId", positionId), zap.String("BaseCurrency", position.BaseCurrency))
		return d.refundUTXO(ctx, u)
	}
	quoteAsset, err := core.GetAssetIdByCurrency(position.QuoteCurrency)
	if err != nil {
		zap.L().Error("exercise:GetAssetIdByCurrency", zap.String("positionId", positionId), zap.String("QuoteCurrency", position.QuoteCurrency))
		return d.refundUTXO(ctx, u)
	}

	if optionType == core.OptionTypePUT {
		if u.AssetID != baseAsset {
			zap.L().Info("exercise:wrong baseAsset", zap.String("positionId", positionId), zap.String("baseAsset", u.AssetID))
			return d.refundUTXO(ctx, u)
		}
		if !(number.FromFloat(position.Size).Equal(number.FromString(size))) ||
			!(number.FromFloat(position.Size).Equal(number.FromFloat(u.Amount))) {
			zap.L().Info("exercise:wrong size or amount", zap.String("positionId", positionId),
				zap.Float64("size", position.Size), zap.Float64("amount", u.Amount))
			return d.refundUTXO(ctx, u)
		}
	}
	if optionType == core.OptionTypeCALL {
		if u.AssetID != quoteAsset {
			zap.L().Info("exercise:wrong quoteAsset", zap.String("positionId", positionId), zap.String("quoteAsset", u.AssetID))
			return d.refundUTXO(ctx, u)
		}
		expectAmount := number.FromFloat(position.Size).Mul(number.FromString(strconv.Itoa(int(instrument.StrikePrice))))
		if !(number.FromFloat(position.Size).Equal(number.FromString(size))) ||
			!(expectAmount.Equal(number.FromFloat(u.Amount))) {
			zap.L().Info("exercise:wrong size or amount", zap.String("positionId", positionId),
				zap.Float64("size", position.Size), zap.Float64("amount", u.Amount))
			return d.refundUTXO(ctx, u)
		}
	}
	if err := d.positionStore.ExercisePositionWithID(ctx, positionId, instrumentName, u.CreatedAt); err != nil {
		zap.L().Error("positionStore.positionExercised", zap.Error(err))
		return err
	}
	return nil
}

// exercise Handler
func (d *UtxoDispatcher) autoExerciseHandler(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {

	zap.L().Info("auto exercised request", zap.String("sendId", u.SenderID))
	if decimal.NewFromFloat(u.Amount).Cmp(gasAmount) != 0 {
		return nil
	}

	var (
		dateString = action.T
		price      string
		err        error
	)

	if price, err = d.deliveryPriceStore.ReadPrice(ctx, core.BTC, dateString); err != nil {
		return err
	}

	deliveryPrice, err := decimal.NewFromString(price)
	if err != nil {
		zap.L().Info("delivery price invalid ,auto exercise passed", zap.String("price", price))
		return err
	}
	if deliveryPrice.LessThanOrEqual(decimal.Zero) {
		zap.L().Info("delivery price invalid ,auto exercise passed", zap.String("price", price))
		return nil
	}

	if err = d.positionService.AutoExercise(ctx, dateString, deliveryPrice, u.CreatedAt); err != nil {
		return err
	}
	return nil
}

// cancelOrderHandler
func (d *UtxoDispatcher) cancelOrderHandler(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {

	orderID := action.O

	//system auto cancel order task
	if action.T == "SYSTEM" && decimal.NewFromFloat(u.Amount).Cmp(gasAmount) == 0 {
		zap.L().Info("system batch cancel order", zap.String("orderId", orderID))
		//todo check time
		o, err := d.orderStore.FindByOrderId(ctx, orderID)
		if err != nil {
			return err
		}
		if o.ID == 0 {
			zap.L().Info("system batch cancel order,order not founded", zap.String("orderId", orderID))
			return nil
		}
		if o.ID > 0 {
			err := d.orderService.CancelOrderAction(ctx, orderID, u.CreatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	}
	//user cancel order via option dance webapp
	if decimal.NewFromFloat(u.Amount).Cmp(gasAmount) == 0 {
		err := d.orderService.CancelOrderAction(ctx, orderID, u.CreatedAt)
		if err != nil {
			return err
		}
	}
	zap.L().Info("cancel order finished", zap.String("orderId", orderID), zap.String("type", action.T))
	return nil
}

func (d *UtxoDispatcher) refundUTXO(ctx context.Context, s *core.UTXO) error {
	return d.transferService.CreateRefundTransfer(ctx, s)
}

func (d *UtxoDispatcher) ExpiringCancelOrder(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {
	zap.L().Info("ProcessUtxo:ExpiringCancelOrder", zap.String("senderId", u.SenderID))
	if decimal.NewFromFloat(u.Amount).Cmp(gasAmount) != 0 {
		return nil
	}
	zap.L().Debug("ProcessUtxo:ExpiringCancelOrder finished")
	return nil
}

func (d *UtxoDispatcher) ExpiryNotify(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {
	zap.L().Info("ProcessUtxo:ExpiryNotify", zap.String("senderId", u.SenderID))
	if decimal.NewFromFloat(u.Amount).Cmp(gasAmount) != 0 {
		return nil
	}
	dayType := action.T
	err := d.notifier.ExpiryNotify(ctx, dayType, u.CreatedAt)
	if err != nil {
		return err
	}
	zap.L().Debug("ProcessUtxo:ExpiryNotify finished", zap.String("dayType", dayType))
	return nil
}

func (d *UtxoDispatcher) Settlement(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {
	zap.L().Info("ProcessUtxo:Settlement", zap.String("senderId", u.SenderID), zap.String("date", action.T), zap.String("type", action.I))
	if decimal.NewFromFloat(u.Amount).Cmp(gasAmount) != 0 {
		return nil
	}
	//check time
	err := d.positionService.Settlement(ctx, "", action.T, action.I, u.CreatedAt)
	if err != nil {
		return err
	}
	zap.L().Debug("ProcessUtxo:Settlement finished")
	return nil
}

func (d *UtxoDispatcher) SyncDbDeliveryPrice(ctx context.Context, u *core.UTXO, action *core.MsgPack) error {
	var (
		date         = action.T
		price        = action.P
		baseCurrency = action.BC
	)
	zap.L().Info("start SyncDbDeliveryPrice", zap.String("date", date), zap.String("price", price))
	_, err := time.Parse(time2.RFC3339Date, date)
	if err != nil {
		zap.L().Error("SyncDbDeliveryPrice wrong date format", zap.String("date", price))
		return d.refundUTXO(ctx, u)
	}
	_, err = decimal.NewFromString(price)
	if err != nil {
		zap.L().Error("SyncDbDeliveryPrice wrong price format", zap.String("price", price))
		return d.refundUTXO(ctx, u)
	}
	err = d.deliveryPriceStore.WritePrice(ctx, baseCurrency, date, price, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

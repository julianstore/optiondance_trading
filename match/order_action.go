package match

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"go.uber.org/zap"
	"option-dance/cmd/config"
	"option-dance/core"
	"time"
)

const (
	pricePrecision  = 4
	AmountPrecision = 4
	fundsPrecision  = 8
)

func (ex *Exchange) PollOrderActions(ctx context.Context) error {

	var (
		checkpoint = 0
		err        error
	)
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			checkpoint, err = ex.LoopPollOrderActions(ctx, checkpoint)
			if err == nil {
				dur = 100 * time.Millisecond
			} else {
				zap.L().Info("PollOrderActions error: ", zap.Error(err))
				dur = 300 * time.Millisecond
			}
		}
	}
}

func (ex *Exchange) LoopPollOrderActions(ctx context.Context, checkpoint int) (res int, err error) {

	var limit = 500
	actions, err := ex.orderService.ListPendingActions(ctx, checkpoint, limit)
	if err != nil {
		return checkpoint, err
	}
	for _, a := range actions {
		ex.ensureProcessOrderAction(ctx, a)
		checkpoint = int(a.ID)
	}
	return checkpoint, nil
}

func (ex *Exchange) ensureProcessOrderAction(ctx context.Context, action *core.Action) {
	o := action.Order
	market := o.InstrumentName

	ex.booksMutex.Lock()
	book := ex.books[market]
	if book == nil {
		book = ex.buildBook(ctx, market)
		go book.Run(ctx)
		ex.books[market] = book
	}
	ex.booksMutex.Unlock()

	price := number.FromString(o.Price).Integer(pricePrecision)
	remainingAmount := number.FromString(o.RemainingAmount).Integer(AmountPrecision)
	filledAmount := number.FromString(o.FilledAmount).Integer(AmountPrecision)
	remainingFunds := number.FromString(o.RemainingFunds).Integer(fundsPrecision)
	filledFunds := number.FromString(o.FilledFunds).Integer(fundsPrecision)
	ex.AttachOrderEvent(ctx, &core.EngineOrder{
		Id:              o.OrderID,
		Side:            o.Side,
		Type:            o.OrderType,
		Price:           price,
		RemainingAmount: remainingAmount,
		FilledAmount:    filledAmount,
		RemainingFunds:  remainingFunds,
		FilledFunds:     filledFunds,
		Quote:           o.QuoteAssetID,
		Base:            o.BaseAssetID,
		BaseCurrency:    o.BaseCurrency,
		QuoteCurrency:   o.QuoteCurrency,
		UserId:          o.UserID,
		BrokerId:        config.Cfg.DApp.AppID,
		InstrumentName:  o.InstrumentName,
		OptionType:      o.OptionType,
		Margin:          o.Margin,
		CreatedAt:       o.CreatedAt,
	}, action)
}

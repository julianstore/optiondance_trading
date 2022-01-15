package mtg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/pkg/mixin"
	"option-dance/pkg/util"
	"option-dance/pkg/util/lock"
	"time"
)

const (
	PageSideASK = "ASK"
	PageSideBID = "BID"
)

var keyLock lock.KeyLock

func CreateOrderRequest(ctx context.Context, action core.OrderAction) (d core.MultisigResponseDTO, err error) {
	//check time
	now := time.Now().Unix()
	instrument, err := core.ParseInstrument(action.InstrumentName)
	if err != nil {
		return d, err
	}
	expiryTs := instrument.ExpirationTimestamp
	if core.Expired(now, expiryTs) {
		return d, fmt.Errorf("cannot order an expired instrument")
	}

	var (
		margin, funds            decimal.Decimal
		memoSide, transferAmount string
	)
	//calculate margin and funds
	if action.Side == core.PageSideAsk {
		memoSide = "A"
		if action.OptionType == core.OptionTypePUT {
			margin = decimal.NewFromInt(instrument.StrikePrice).Mul(decimal.RequireFromString(action.Amount))
		} else if action.OptionType == core.OptionTypeCALL {
			margin = decimal.RequireFromString(action.Amount)
		} else {
			return d, fmt.Errorf("wrong option Type %s", action.OptionType)
		}
		transferAmount = margin.String()
	} else if action.Side == core.PageSideBid {
		memoSide = "B"
		funds = decimal.RequireFromString(action.Price).Mul(decimal.RequireFromString(action.Amount))
		transferAmount = funds.String()
	} else {
		return d, fmt.Errorf("wrong option side %s", action.Side)
	}
	orderMsg := core.MsgPack{
		S:  memoSide,
		T:  action.Type,
		I:  action.InstrumentName,
		P:  action.Price,
		A:  action.Amount,
		TC: action.TraceId,
		Q:  action.QuoteAsset,
		B:  action.BaseAsset,
		QC: action.QuoteCurrency,
		BC: action.BaseCurrency,
		M:  margin.String(),
		AC: core.MsgPackActionCreateOrder,
	}
	// check msg pack
	if err = orderMsg.CheckCreateOrderPack(); err != nil {
		return d, err
	}
	memo, err := orderMsg.Pack()
	if err != nil {
		return d, err
	}

	var (
		dapp             = config.Cfg.DApp
		transferAsset    string
		transferCurrency string
	)
	//asset id
	if action.Side == core.PageSideBid {
		transferCurrency = action.QuoteCurrency
	} else {
		if action.OptionType == core.OptionTypePUT {
			transferCurrency = action.QuoteCurrency
		} else {
			transferCurrency = action.BaseCurrency
		}
	}
	transferAsset, err = core.GetAssetIdByCurrency(transferCurrency)
	if err != nil {
		return d, err
	}
	multisigDTO := core.MultisigDTO{
		AssetId: transferAsset,
		Amount:  transferAmount,
		TraceId: action.TraceId,
		Memo:    memo,
		OpponentMultisig: core.OpponentMultisig{
			Receivers: dapp.Receivers,
			Threshold: dapp.Threshold,
		},
	}
	marshal, err := json.Marshal(multisigDTO)
	if err != nil {
		return d, err
	}
	res, err := mixin.MixinRequest("/payments", "POST", marshal)
	if err != nil {
		return d, err
	}
	if err = json.Unmarshal(res, &d); err != nil {
		return d, err
	}
	return d, nil
}

//transfer 0.00000001 CNB with orderId msgPack in memo to multisig to cancel order
func CancelOrderRequest(ctx context.Context, orderId, userId string) error {
	if !keyLock.TryLock(orderId) {
		return errors.New("requests are too frequent")
	}
	defer keyLock.UnLock(orderId)
	var o *core.Order
	if err := config.Db.Model(core.Order{}).Where(" user_id = ? AND order_id = ?", userId, orderId).First(&o).Error; err != nil {
		return err
	}
	if o.OrderStatus != core.OrderStatusOpen {
		return errors.New("non-open order cannot be cancelled")
	}
	dapp := config.Cfg.DApp
	// msg pack
	cnbAmount := "0.00000001"
	orderMsg := core.MsgPack{
		O:  orderId,
		AC: core.MsgPackActionCancelOrder,
	}
	memo, err := orderMsg.Pack()
	if err != nil {
		return err
	}
	cnbAssetID, _ := core.GetAssetIdByCurrency("CNB")
	input := bot.TransferInput{
		AssetId: cnbAssetID,
		Amount:  number.FromString(cnbAmount),
		TraceId: util.UuidNewV4().String(),
		Memo:    memo,
		OpponentMultisig: struct {
			Receivers []string
			Threshold int64
		}{
			Receivers: dapp.Receivers,
			Threshold: dapp.Threshold,
		},
	}

	_, err = bot.CreateMultisigTransaction(ctx, &input, dapp.AppID, dapp.SessionID, dapp.PrivateKey, dapp.Pin, dapp.PinToken)
	if err != nil {
		return err
	}

	if err := config.Db.Model(core.Order{}).Where("order_id = ? and order_status = ?", orderId, core.OrderStatusOpen).
		Updates(map[string]interface{}{
			"order_status": core.OrderStatusCanceling,
			"updated_at":   time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}

//exercise
func ExerciseRequest(ctx context.Context, instrumentName, optionSize, positionId string) (d *core.MultisigResponseDTO, err error) {
	instrument, err := core.ParseInstrument(instrumentName)
	if err != nil {
		return nil, err
	}
	exerciseTs := time.Now().Unix()
	expiryTs := instrument.ExpirationTimestamp
	if !core.Exercisable(exerciseTs, expiryTs) {
		return nil, fmt.Errorf("exercise at wrong time")
	}

	assetId, err := core.GetAssetIdByCurrency(instrument.BaseCurrency)
	if err != nil {
		return nil, err
	}
	dapp := config.Cfg.DApp
	exerciseMsg := core.MsgPack{
		AC: core.MsgPackActionExercise,
		I:  instrumentName,
		A:  optionSize,
		TC: positionId,
	}
	memo, err := exerciseMsg.Pack()
	if err != nil {
		return nil, err
	}
	multisigDTO := core.MultisigDTO{
		AssetId: assetId,
		Amount:  optionSize,
		TraceId: util.UuidNewV4().String(),
		Memo:    memo,
		OpponentMultisig: core.OpponentMultisig{
			Receivers: dapp.Receivers,
			Threshold: dapp.Threshold,
		},
	}
	marshal, err := json.Marshal(multisigDTO)
	if err != nil {
		return nil, err
	}
	res, err := mixin.MixinRequest("/payments", "POST", marshal)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func CashDeliveryExerciseTx(ctx context.Context, instrumentName, optionSize, positionId string) (*bot.RawTransaction, error) {
	pack := core.MsgPack{
		AC: core.MsgPackActionExercise,
		I:  instrumentName,
		A:  optionSize,
		TC: positionId,
	}
	return CreateMultiSignTx(ctx, pack)
}

func SettlementTx(ctx context.Context, date, deliveryType string) (*bot.RawTransaction, error) {
	pack := core.MsgPack{
		T:  date,
		AC: "SETTLEMENT",
		I:  deliveryType,
	}
	return CreateMultiSignTx(ctx, pack)
}

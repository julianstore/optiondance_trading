package notifier

import (
	"context"
	"crypto/md5"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"io"
	"math/big"
	"option-dance/core"
	"option-dance/service/message"
	"strconv"
	"time"
)

type notifier struct {
	messageStore   core.MessageStore
	positionStore  core.PositionStore
	messageBuilder core.MessageBuilder
}

func NewNotifier(
	messageStore core.MessageStore,
	positionStore core.PositionStore,
	messageBuilder core.MessageBuilder,
) core.Notifier {
	return &notifier{
		messageStore:   messageStore,
		positionStore:  positionStore,
		messageBuilder: messageBuilder,
	}
}

func (s *notifier) SnapshotCard(ctx context.Context, transfer *core.Transfer, signedTx string) error {
	tx, err := mixin.TransactionFromRaw(signedTx)
	if err != nil {
		zap.L().Error("TransactionFromRaw", zap.Error(err))
		return nil
	}

	hash, err := tx.TransactionHash()
	if err != nil {
		return nil
	}

	if transfer.Threshold != 1 {
		return nil
	}
	traceID := mixinRawTransactionHashTraceID(hash.String(), 0)
	msg, err := message.BuildTransferCardMsg(ctx, transfer.UserId, transfer.Amount, transfer.AssetId, traceID)
	if err != nil {
		return err
	}
	err = s.messageStore.Create(ctx, &msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *notifier) Transfer(ctx context.Context, transfer *core.Transfer) error {
	if transfer.Source == core.TransferSourceUTXORefund ||
		transfer.Source == core.TransferSourceOrderInvalid ||
		transfer.Source == core.TransferSourcePositionExercisedInvalid ||
		transfer.Source == core.TransferSourceUTXOMerge {
		return nil
	}
	m, err := message.BuildBizInfoMsg(ctx, transfer.Source, transfer.UserId, transfer)
	if err != nil {
		return err
	}
	zap.L().Debug("transfer msg", zap.Any("msg", m))
	return s.messageStore.Create(ctx, &m)
}

func (s *notifier) Trade(ctx context.Context, trade *core.Trade, transfer *core.Transfer) error {
	m, err := message.GetBizInfoMsgByTrade(ctx, core.TransferSourceTradeConfirmed, trade, transfer)
	if err != nil {
		return err
	}
	return s.messageStore.Create(ctx, &m)
}

func (s *notifier) ExpiryNotify(ctx context.Context, day string, notifyTime time.Time) error {
	if day == "day0" {
		list, err := s.positionStore.ListInstrumentNamesByDate(ctx, notifyTime, core.PositionStatusALL)
		if err != nil {
			return err
		}
		err = s.createExpiryNotify(ctx, list, "ntf-day0")
		if err != nil {
			return err
		}
	}
	if day == "day1-2" {
		day1 := notifyTime.Add(time.Hour * 24)
		day2 := notifyTime.Add(time.Hour * 48)
		list1, err := s.positionStore.ListInstrumentNamesByDate(ctx, day1, core.PositionStatusALL)
		if err != nil {
			return err
		}
		err = s.createExpiryNotify(ctx, list1, "ntf-day1")
		if err != nil {
			return err
		}
		list2, err := s.positionStore.ListInstrumentNamesByDate(ctx, day2, core.PositionStatusALL)
		if err != nil {
			return err
		}
		err = s.createExpiryNotify(ctx, list2, "ntf-day2")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *notifier) createExpiryNotify(ctx context.Context, nameList []string, modifier string) error {
	zap.L().Info("Cron-ExpiryNotify", zap.Strings("nameList", nameList))
	for _, name := range nameList {
		positions, err := s.positionStore.ListExercisablePosition(ctx, name)
		if err != nil {
			return err
		}
		for _, p := range positions {
			notifyContent, err := message.GetExpiryNotifyContent(p.InstrumentName, p.Side, p.Size)
			if err != nil {
				return nil
			}
			msg, err := message.BuildExerciseNotifyMsg(ctx, p.UserID, p.PositionID, notifyContent, modifier)
			if err != nil {
				return nil
			}
			err = s.messageStore.Create(ctx, &msg)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func (s *notifier) CreateOrderGroupNotify(ctx context.Context, order core.EngineOrder) error {
	instrument, err := core.ParseInstrument(order.InstrumentName)
	if err != nil {
		return err
	}

	var remainingAmount = order.RemainingAmount.Add(order.FilledAmount).Persist()
	if order.Side == core.PageSideBid {
		remainingAmount = (order.RemainingFunds.Add(order.FilledFunds)).Div(order.Price).Persist()
	}

	o := core.Order{
		OrderID:             order.Id,
		UserID:              order.UserId,
		Side:                order.Side,
		OrderType:           order.Type,
		Price:               order.Price.Persist(),
		RemainingAmount:     remainingAmount,
		FilledAmount:        order.FilledAmount.Persist(),
		RemainingFunds:      order.RemainingFunds.Persist(),
		FilledFunds:         order.FilledFunds.Persist(),
		Margin:              order.Margin,
		InstrumentName:      order.InstrumentName,
		OptionType:          order.OptionType,
		StrikePrice:         strconv.Itoa(int(instrument.StrikePrice)),
		ExpirationTimestamp: instrument.ExpirationTimestamp,
		ExpirationDate:      instrument.ExpirationDate,
		QuoteCurrency:       instrument.QuoteCurrency,
		BaseCurrency:        instrument.BaseCurrency,
	}
	msg, err := s.messageBuilder.OrderCreateGroupNotify(ctx, &o)
	if err != nil {
		return err
	}
	if err = s.messageStore.Create(ctx, msg); err != nil {
		return err
	}
	return nil
}

func mixinRawTransactionHashTraceID(hash string, index uint8) string {
	h := md5.New()
	_, _ = io.WriteString(h, hash)
	b := new(big.Int).SetInt64(int64(index))
	h.Write(b.Bytes())
	s := h.Sum(nil)
	s[6] = (s[6] & 0x0f) | 0x30
	s[8] = (s[8] & 0x3f) | 0x80
	sid, err := uuid.FromBytes(s)
	if err != nil {
		panic(err)
	}

	return sid.String()
}

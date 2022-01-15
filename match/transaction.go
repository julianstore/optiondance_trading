package match

import (
	"context"
	"crypto/md5"
	"errors"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"io"
	"math/big"
	"option-dance/core"
	"time"
)

const (
	SendRawTransactionRetryLimit = 20
)

func (ex *Exchange) LoopingSendTransaction(ctx context.Context) error {

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := ex.LoopSendTransaction(ctx); err == nil {
				dur = 200 * time.Millisecond
			} else {
				zap.L().Info("LoopingSendTransaction error: ", zap.Error(err))
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (ex *Exchange) LoopSendTransaction(ctx context.Context) error {

	limit := 500
	rawTransactions, err := ex.rawTxStore.ListPendingRawTransactions(ctx, limit)
	if err != nil {
		return err
	}
	for _, t := range rawTransactions {
		err := ex.mustSendRawTransaction(ctx, t)
		if err == nil {
			_ = ex.rawTxStore.ExpireRawTransaction(ctx, t)
		}
	}
	return nil
}

func (ex *Exchange) mustSendRawTransaction(ctx context.Context, tx *core.RawTransaction) error {
	retryTimes := 0
	for {
		retryTimes++
		err := ex.sendRawTransaction(ctx, tx.Data)
		if err == nil {
			break
		} else {
			if retryTimes >= SendRawTransactionRetryLimit {
				//tx.State = "invalid"
				//tx.ErrorInfo = err.Error()
				//if err := core.SaveRawTransaction(ctx, tx); err != nil {
				//	zap.L().Error("SaveRawTransactionError", zap.Error(err))
				//}
				//if err = utxo.Invalid(ctx, tx.TraceID); err != nil {
				//	zap.L().Error("utxo.Invalid", zap.Error(err))
				//}
				//return fmt.Errorf("%s", tx.ErrorInfo)
				break
			}
			time.Sleep(PollInterval)
		}
	}
	return nil
}

func (ex *Exchange) sendRawTransaction(ctx context.Context, raw string) error {
	if tx, err := mixin.SendRawTransaction(ctx, raw); err != nil {
		return err
	} else if tx.Snapshot != nil {
		return nil
	}

	var txHash mixin.Hash
	if tx, err := mixin.TransactionFromRaw(raw); err == nil {
		txHash, _ = tx.TransactionHash()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return errors.New("mixin net snapshot not generated")
		case <-time.After(dur):
			if tx, err := mixin.GetTransaction(ctx, txHash); err != nil {
				return err
			} else if tx.Snapshot != nil {
				return nil
			}
			dur = time.Second
		}
	}
}

func RawTransactionTraceID(raw string) (string, error) {
	var txHash mixin.Hash
	if tx, err := mixin.TransactionFromRaw(raw); err == nil {
		txHash, _ = tx.TransactionHash()
	} else {
		return "", err
	}
	id := MixinRawTransactionHashTraceID(txHash.String(), 0)
	return id, nil
}

func MixinRawTransactionHashTraceID(hash string, index uint8) string {
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

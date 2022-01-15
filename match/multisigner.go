package match

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"option-dance/cmd/config"
	"option-dance/core"
	m "option-dance/pkg/mixin"
	"strings"
	"time"

	"github.com/MixinNetwork/go-number"
	"github.com/asaskevich/govalidator"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/ugorji/go/codec"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TransferAction struct {
	S string // source
	O string // cancelled order
	A string // matched ask order
	B string // matched bid order
}

type MultiSigner struct {
	codec         codec.Handle
	notifier      core.Notifier
	members       []string
	threshold     uint8
	utxoStore     core.UtxoStore
	transferStore core.TransferStore
	rawTxStore    core.RawTxStore
}

func NewMultiSigner(
	notifier core.Notifier,
	utxoStore core.UtxoStore,
	rawTxStore core.RawTxStore,
	transferStore core.TransferStore,
) MultiSigner {
	return MultiSigner{
		codec:         new(codec.MsgpackHandle),
		notifier:      notifier,
		members:       config.Cfg.DApp.Receivers,
		threshold:     uint8(config.Cfg.DApp.Threshold),
		utxoStore:     utxoStore,
		transferStore: transferStore,
		rawTxStore:    rawTxStore,
	}
}

func (s *MultiSigner) Run(ctx context.Context) error {
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := s.handlePendingTransfers(ctx); err == nil {
				dur = 200 * time.Millisecond
			} else {
				zap.L().Info("LoopingSendMessage error: ", zap.Error(err))
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (s *MultiSigner) handlePendingTransfers(ctx context.Context) error {

	limit := 500
	transfers, err := s.transferStore.ListPendingTransfers(ctx, limit)
	if err != nil {
		return err
	}
	for _, t := range transfers {
		if err := s.processTransfer(ctx, t); err != nil {
			return err
		}
	}
	if len(transfers) < limit {
		time.Sleep(PollInterval)
	}
	return nil
}

func (s *MultiSigner) processTransfer(ctx context.Context, transfer *core.Transfer) error {
	if number.FromString(transfer.Amount).Exhausted() {
		zap.L().Info("processTransfer skipped", zap.String("transferID", transfer.TransferId))
		return nil
	}
	transfer.OpponentList = strings.Split(transfer.Opponents, ",")
	var data *TransferAction
	switch transfer.Source {
	case core.TransferSourceOrderFilled:
		data = &TransferAction{S: "FILL", O: transfer.Detail}
	case core.TransferSourceOrderCancelled:
		data = &TransferAction{S: "CANCEL", O: transfer.Detail}
	case core.TransferSourceOrderInvalid:
		data = &TransferAction{S: "REFUND", O: transfer.Detail}
	case core.TransferSourceTradeConfirmed:
		trade, err := s.transferStore.ReadTransferTrade(ctx, transfer.Detail, transfer.AssetId)
		if err != nil {
			return err
		}
		data = &TransferAction{S: "MATCH", A: "", B: ""}
		if trade != nil {
			data.A = trade.AskOrderID
			data.B = trade.BidOrderID
		} else {
			zap.L().Warn("processTransfer warn:trade nil", zap.Any("transfer", transfer))
		}
	case core.TransferSourcePositionExercised:
		data = &TransferAction{S: core.TransferSourcePositionExercised, O: transfer.Detail}
	case core.TransferSourcePositionExercisedRefund:
		data = &TransferAction{S: core.TransferSourcePositionExercisedRefund, O: transfer.Detail}
	case core.TransferSourcePositionExercisedInvalid:
		data = &TransferAction{S: core.TransferSourcePositionExercisedInvalid, O: transfer.Detail}
	default:
		data = &TransferAction{S: transfer.Source, O: transfer.Detail}
	}
	out := make([]byte, 200)
	encoder := codec.NewEncoderBytes(&out, s.codec)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	memo := base64.StdEncoding.EncodeToString(out)
	if len(memo) > 200 {
		zap.L().Warn("memo length warn", zap.String("memo", memo), zap.Any("transfer", transfer))
	}
	transfer.Memo = memo
	return s.handleTransferV2(ctx, transfer)
}

func (s *MultiSigner) handleTransferV2(ctx context.Context, transfer *core.Transfer) error {

	const limit = 32
	outputs, err := s.utxoStore.ListUnspentV2(ctx, transfer.AssetId, limit)
	if err != nil {
		zap.L().Error("utxo:ListUnspent", zap.Error(err))
		return err
	}

	var (
		idx         int
		sum         decimal.Decimal
		traces      []string
		totalAmount = decimal.RequireFromString(transfer.Amount)
		dApp        = config.Cfg.DApp
	)

	for _, output := range outputs {
		sum = sum.Add(decimal.NewFromFloat(output.Amount))
		traces = append(traces, output.TraceID)
		idx++

		if sum.GreaterThanOrEqual(totalAmount) {
			break
		}
	}

	outputs = outputs[:idx]

	if sum.LessThan(totalAmount) {
		// merge outputs
		if len(outputs) == limit {
			traceID := Modify(transfer.TraceID, mixin.HashMembers(traces))
			merge := &core.Transfer{
				TransferId: traceID,
				Source:     core.TransferSourceUTXOMerge,
				TraceID:    traceID,
				AssetId:    transfer.AssetId,
				Amount:     sum.String(),
				Opponents:  strings.Join(dApp.Receivers, ","),
				Threshold:  int8(dApp.Threshold),
				Memo:       fmt.Sprintf("merge for %s", transfer.TraceID),
			}
			return s.spentV2(ctx, outputs, merge)
		}

		err := errors.New("insufficient balance")
		zap.L().Error("handle transfer error", zap.Error(err))
		return err
	}
	return s.spentV2(ctx, outputs, transfer)
}

func (s *MultiSigner) spentV2(ctx context.Context, outputs []*core.UTXO, transfer *core.Transfer) error {
	if tx, err := s.Spend(ctx, outputs, transfer); err != nil {
		return err
	} else if tx != nil {
		tx.UserID = transfer.UserId
		tx.AssetID = transfer.AssetId
		tx.Amount = transfer.Amount
		if err = s.rawTxStore.SaveRawTransaction(ctx, tx); err != nil {
			zap.L().Error("SaveRawTransaction Error", zap.Error(err))
			return err
		}
	}
	if err := s.utxoSpentBy(outputs, transfer); err != nil {
		zap.L().Error("UTXO Spent Error", zap.Error(err))
		return err
	}
	err := s.notifier.Transfer(ctx, transfer)
	if err != nil {
		zap.L().Error("notifier.Transfer", zap.Error(err))
	}
	return nil
}

func Modify(id, modifier string) string {
	ns, err := uuid.FromString(id)
	if err != nil {
		zap.L().Error("uuidModify", zap.Error(err))
	}
	return uuid.NewV5(ns, modifier).String()
}

func (s *MultiSigner) utxoSpentBy(utxos []*core.UTXO, transfer *core.Transfer) error {
	if err := config.Db.Transaction(func(tx *gorm.DB) error {
		for _, output := range utxos {
			if err := tx.Model(core.UTXO{}).Where("trace_id=?", output.TraceID).Updates(map[string]interface{}{
				"spent_by": transfer.TraceID,
			}).Error; err != nil {
				return err
			}
		}
		if transfer.Source == core.TransferSourceUTXOMerge {
			transfer.Handled = 1
			if err := tx.Create(transfer).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(transfer).Where("transfer_id = ?", transfer.TransferId).Updates(map[string]interface{}{
				"handled": 1,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Spent Consumption specified UTXO
// If transfer is nil, merge these UTXOs
func (s *MultiSigner) Spend(ctx context.Context, outputs []*core.UTXO, transfer *core.Transfer) (*core.RawTransaction, error) {
	state, tx, err := s.signTransaction(ctx, outputs, transfer)
	if err != nil {
		return nil, err
	}

	switch state {
	case mixin.UTXOStateSpent:
		//  m There are n node signatures in the node, and Output has all been spent. Here is a simple verification:
		//	1. memo is as expected
		//	2. The amount is the same as expected (the first output is expenditure, and the second output is change)
		//	3. The first output is a normal transfer and is transferred to a single address
		//	4. If there is a second output, it should also be a multi-signature of n/m
		tx, err := mixin.TransactionFromRaw(tx)
		if err != nil {
			return nil, err
		}

		if err := s.validateTransaction(tx, transfer); err != nil {
			return nil, fmt.Errorf("validateTransaction failed: %w", err)
		}

	case mixin.UTXOStateSigned:
		client, err := m.Client()
		if err != nil {
			return nil, err
		}
		sig, err := client.CreateMultisig(ctx, mixin.MultisigActionSign, tx)
		if err != nil {
			return nil, fmt.Errorf("CreateMultisig %s failed: %w", mixin.MultisigActionSign, err)
		}
		app := config.Cfg.DApp
		if !govalidator.IsIn(app.AppID, sig.Signers...) {
			if valiErr := s.validateMultisig(sig, transfer); valiErr != nil {
				return nil, valiErr
			}

			sig, err = client.SignMultisig(ctx, sig.RequestID, app.Pin)
			if err != nil {
				return nil, fmt.Errorf("SignMultisig failed: %w", err)
			}
		}
		if len(sig.Signers) >= int(sig.Threshold) {
			return &core.RawTransaction{
				TraceID: transfer.TraceID,
				Data:    sig.RawTransaction,
			}, nil
		}
	default:
		return nil, errors.New("cannot consume unsigned utxo")
	}
	return nil, nil
}

func (s *MultiSigner) signTransaction(ctx context.Context, outputs []*core.UTXO, transfer *core.Transfer) (string, string, error) {
	if len(outputs) == 0 {
		return mixin.UTXOStateSpent, "", nil
	}
	input := &mixin.TransactionInput{
		Memo: transfer.Memo,
		Hint: transfer.TransferId,
	}

	state := outputs[0].State
	signedTx := outputs[0].SignedTx
	sum := decimal.Zero

	for _, output := range outputs[0:] {
		st := output.State
		tx := output.SignedTx
		sum = sum.Add(output.Raw.Amount)

		if st == state && tx == signedTx {
			input.AppendUTXO(output.Raw)
			continue
		}

		return "", "", errors.New("state not match")
	}

	if signedTx != "" {
		return state, signedTx, nil
	}

	receivers := []string{transfer.UserId}
	if len(transfer.UserId) == 0 {
		receivers = strings.Split(transfer.Opponents, ",")
	}

	input.AppendOutput(receivers, uint8(transfer.Threshold), decimal.RequireFromString(transfer.Amount))
	client, err2 := m.Client()
	if err2 != nil {
		return "", "", err2
	}
	tx, err := client.MakeMultisigTransaction(ctx, input)
	if err != nil {
		return "", "", err
	}
	signedTx, _ = tx.DumpTransaction()
	return mixin.UTXOStateSigned, signedTx, nil
}

// validateTransaction validate spent Tx
func (s *MultiSigner) validateTransaction(tx *mixin.Transaction, transfer *core.Transfer) error {
	if string(tx.Extra) != transfer.Memo {
		return fmt.Errorf("memo not match, expect %q got %q", transfer.Memo, string(tx.Extra))
	}

	for idx, output := range tx.Outputs {
		switch idx {
		case 0: // Check output and transfer
			if output.Type != 0 {
				return fmt.Errorf("first output type not matched, expect %d got %d", 0, output.Type)
			}

			if expect, got := mixin.NewIntegerFromString(transfer.Amount).String(), output.Amount.String(); expect != got {
				return fmt.Errorf("amount not match, expect %s got %s", expect, got)
			}

			if expect, got := mixin.NewThresholdScript(uint8(transfer.Threshold)).String(), output.Script.String(); expect != got {
				return fmt.Errorf("first output script not matched, expect %s got %s", expect, got)
			}

			opponents := strings.Split(transfer.Opponents, ",")
			if len(output.Keys) != len(opponents) {
				return errors.New("receivers not match")
			}

		default: // Check change
			if expect, got := mixin.NewThresholdScript(s.threshold).String(), output.Script.String(); expect != got {
				return fmt.Errorf("first output script not matched, expect %s got %s", expect, got)
			}

			if len(output.Keys) != len(s.members) {
				return errors.New("receivers not match")
			}
		}
	}
	return nil
}

// validateMultisig validate multisig request
func (s *MultiSigner) validateMultisig(req *mixin.MultisigRequest, transfer *core.Transfer) error {
	if req.AssetID != transfer.AssetId {
		return fmt.Errorf("asset id not match, expect %q got %q", transfer.AssetId, req.AssetID)
	}

	if req.Memo != transfer.Memo {
		return fmt.Errorf("memo not match, expect %q got %q", transfer.Memo, req.Memo)
	}
	fromString, err := decimal.NewFromString(transfer.Amount)
	if err != nil {
		return err
	}
	if !req.Amount.Equal(fromString) {
		return fmt.Errorf("amount not match, expect %s got %s", transfer.Amount, req.Amount)
	}

	opponents := strings.Split(transfer.Opponents, ",")
	if mixin.HashMembers(req.Receivers) != mixin.HashMembers(opponents) {
		return errors.New("receivers not match")
	}

	return nil
}

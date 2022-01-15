package match

import (
	"context"
	"github.com/fox-one/mixin-sdk-go"
	"go.uber.org/zap"
	"option-dance/core"
	"time"
)

type SpentSyncer struct {
	propertyStore core.PropertyStore
	positionStore core.PositionStore
	notifier      core.Notifier
	utxoStore     core.UtxoStore
	transferStore core.TransferStore
}

func NewSpentSyncer(
	propertyStore core.PropertyStore,
	positionStore core.PositionStore,
	utxoStore core.UtxoStore,
	transferStore core.TransferStore,
	notifier core.Notifier,
) SpentSyncer {
	return SpentSyncer{
		propertyStore: propertyStore,
		positionStore: positionStore,
		notifier:      notifier,
		utxoStore:     utxoStore,
		transferStore: transferStore,
	}
}

func (s *SpentSyncer) Run(ctx context.Context) error {

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := s.spentSync(ctx); err == nil {
				dur = 500 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (s *SpentSyncer) spentSync(ctx context.Context) error {

	const limit = 100
	transfers, err := s.transferStore.ListTransfers(ctx, core.TransferStatusHandled, limit)
	if err != nil {
		zap.L().Error("transferStore.ListTransfers", zap.Error(err))
		return err
	}

	if len(transfers) == 0 {
		return nil
	}

	for _, transfer := range transfers {
		_ = s.handleSpentSync(ctx, transfer)
	}

	return nil
}

func (s *SpentSyncer) handleSpentSync(ctx context.Context, transfer *core.Transfer) error {

	output, err := s.utxoStore.FindSpentBy(ctx, transfer.AssetId, transfer.TraceID)
	if err != nil {
		zap.L().Error("utxoStore.FindSpentBy", zap.Error(err))
		return err
	}

	if output.State != mixin.UTXOStateSpent {
		return nil
	}

	signedTx := output.SignedTx
	if signedTx == "" {
		return nil
	}

	if err = s.notifier.SnapshotCard(ctx, transfer, signedTx); err != nil {
		zap.L().Error("notifier.SnapshotCard", zap.Error(err))
		return err
	}

	if err := s.transferStore.UpdateTransfer(ctx, transfer, core.TransferStatusDone); err != nil {
		zap.L().Error("transferStore.UpdateTransfer", zap.Error(err))
		return err
	}
	return nil
}

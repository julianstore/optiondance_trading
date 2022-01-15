package match

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/pkg/mixin"
	"time"
)

const (
	MtgUtxoSyncOffsetKey = "utxo_sync_checkpoint"
)

type UtxoSyncer struct {
	propertyStore core.PropertyStore
	utxoStore     core.UtxoStore
}

func NewUtxoSyncer(propertyStore core.PropertyStore, utxoStore core.UtxoStore) UtxoSyncer {
	s := UtxoSyncer{
		propertyStore: propertyStore,
		utxoStore:     utxoStore,
	}
	return s
}

func (s *UtxoSyncer) Run(ctx context.Context) error {
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := s.SyncUTXO(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				zap.L().Info("SyncUtxo error: ", zap.Error(err))
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (s *UtxoSyncer) SyncUTXO(ctx context.Context) error {
	RFC3339NanoString, err := s.propertyStore.ReadProperty(ctx, MtgUtxoSyncOffsetKey)
	if err != nil {
		return err
	}
	offset, err := time.Parse(time.RFC3339Nano, RFC3339NanoString)
	if err != nil {
		err = s.propertyStore.WriteTimeProperty(ctx, MtgUtxoSyncOffsetKey, time.Unix(0, 0))
		if err != nil {
			return err
		}
	}

	var (
		outputs   = make([]core.UTXO, 0, 8)
		positions = make(map[string]int)
		pos       = 0
		limit     = 500
	)

	for {
		utxos, err := s.PullLatestUTXOs(ctx, limit, offset)
		if err != nil {
			return err
		}
		for _, u := range utxos {
			offset = u.UpdatedAt
			p, ok := positions[u.TraceID]
			if ok {
				outputs[p] = u
				continue
			}
			outputs = append(outputs, u)
			positions[u.TraceID] = pos
			pos++
		}
		if len(utxos) < limit {
			break
		}
	}
	if len(outputs) == 0 {
		return nil
	}
	for _, output := range outputs {
		err := s.utxoStore.Save(ctx, &output)
		if err != nil {
			return err
		}
		err = s.propertyStore.WriteTimeProperty(ctx, MtgUtxoSyncOffsetKey, offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UtxoSyncer) PullLatestUTXOs(ctx context.Context, limit int, offset time.Time) (utxoList []core.UTXO, err error) {

	client, err := mixin.Client()
	if err != nil {
		return nil, err
	}

	dApp := config.Cfg.DApp
	outputs, err := client.ReadMultisigOutputs(ctx, dApp.Receivers, uint8(dApp.Threshold), offset, limit)
	if err != nil {
		return nil, err
	}
	for _, output := range outputs {
		jsonData, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		utxo := core.UTXO{
			UserID:          output.UserID,
			TraceID:         output.UTXOID,
			AssetID:         output.AssetID,
			Memo:            output.Memo,
			SenderID:        output.Sender,
			Amount:          output.Amount.InexactFloat64(),
			TransactionHash: output.TransactionHash.String(),
			State:           output.State,
			SignedBy:        output.SignedBy,
			SignedTx:        output.SignedTx,
			CreatedAt:       output.CreatedAt,
			UpdatedAt:       output.CreatedAt,
			JSONData:        string(jsonData),
			Invalid:         0,
		}
		utxoList = append(utxoList, utxo)
	}
	core.SortUTXOs(utxoList)
	return
}

package notifier

import (
	"context"
	"option-dance/core"
	"time"
)

func NewMuteNotifier() core.Notifier {
	return &muteNotifier{}
}

type muteNotifier struct{}

func (m *muteNotifier) SnapshotCard(ctx context.Context, transfer *core.Transfer, TxHash string) error {
	return nil
}

func (m *muteNotifier) ExpiryNotify(ctx context.Context, day string, notifyTime time.Time) error {
	return nil
}

func (m *muteNotifier) Transfer(ctx context.Context, transfer *core.Transfer) error {
	return nil
}

func (m *muteNotifier) Trade(ctx context.Context, trade *core.Trade, transfer *core.Transfer) error {
	return nil
}

func (s *muteNotifier) CreateOrderGroupNotify(ctx context.Context, order core.EngineOrder) error {
	return nil
}

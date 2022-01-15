package core

import (
	"context"
	"time"
)

type Notifier interface {
	SnapshotCard(ctx context.Context, transfer *Transfer, TxHash string) error
	ExpiryNotify(ctx context.Context, day string, time time.Time) error
	Transfer(ctx context.Context, transfer *Transfer) error
	Trade(ctx context.Context, trade *Trade, transfer *Transfer) error
	CreateOrderGroupNotify(ctx context.Context, order EngineOrder) error
}

package core

import (
	"context"
	"crypto/md5"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"io"
	"time"
)

const (
	TransferSourceUTXORefund               = "UTXO_REFUND"
	TransferSourceUTXOMerge                = "UTXO_MERGE"
	TransferSourceTradeConfirmed           = "TRADE_CONFIRMED"
	TransferSourceTradeBidRefund           = "TRADE_BID_REFUND"
	TransferSourceOrderCancelled           = "ORDER_CANCELLED"
	TransferSourceOrderFilled              = "ORDER_FILLED"
	TransferSourceOrderInvalid             = "ORDER_INVALID"
	TransferSourcePositionExercised        = "POSITION_EXERCISED"
	TransferSourcePositionExercisedRefund  = "POSITION_EXERCISED_REFUND"
	TransferSourcePositionExercisedInvalid = "POSITION_EXERCISED_INVALID"
	TransferSourcePositionClosedRefund     = "POSITION_CLOSED_REFUND"

	TransferStatusPending = 0
	TransferStatusHandled = 1
	TransferStatusDone    = 2
)

type (
	Transfer struct {
		ID           int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		TransferId   string    `gorm:"column:transfer_id;type:varchar(60);not null;unique" json:"transfer_id"`
		Source       string    `gorm:"column:source;type:varchar(45)" json:"source"`
		Detail       string    `gorm:"column:detail;type:varchar(100)" json:"detail"`
		AssetId      string    `gorm:"column:asset_id;type:varchar(45)" json:"asset_id"`
		Amount       string    `gorm:"column:amount;type:varchar(45)" json:"amount"`
		CreatedAt    time.Time `gorm:"column:created_at;type:datetime(6);index:idx_created_at" json:"created_at"`
		UserId       string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		ClientId     string    `gorm:"column:client_id;type:varchar(45)" json:"client_id"`
		Memo         string    `gorm:"column:memo;type:varchar(500)" json:"memo"`
		TraceID      string    `gorm:"column:trace_id;type:varchar(45)" json:"trace_id"`
		Passed       int8      `gorm:"column:passed;type:tinyint" json:"passed"`
		Handled      int8      `gorm:"column:handled;type:tinyint" json:"handled"`
		Opponents    string    `gorm:"column:opponents;type:varchar(1000)" json:"opponents"`
		Threshold    int8      `gorm:"column:threshold;type:tinyint" json:"threshold"`
		OpponentList []string  `gorm:"-" json:"-"`
	}

	TransferStore interface {
		ListTransfers(ctx context.Context, status int8, limit int) ([]*Transfer, error)
		UpdateTransfer(ctx context.Context, transfer *Transfer, status int8) error
		ReadTransferTrade(ctx context.Context, tradeId, AssetId string) (trade *Trade, err error)
		ListPendingTransfers(ctx context.Context, limit int) (transfers []*Transfer, err error)
		Create(ctx context.Context, transfer Transfer) error
		Exist(ctx context.Context, transferId string) (bool, error)
	}

	TransferService interface {
		CreateRefundTransfer(ctx context.Context, s *UTXO) error
		CreateClosePositionTransfer(tx *gorm.DB, margin float64, positionId, userId, optionType string, trade *Trade) error
	}
)

func (m *Transfer) TableName() string {
	return "transfer"
}

func MemoToAction(ctx context.Context, u *UTXO) (*MsgPack, error) {
	pack, err := DecryptAction(u.Memo)
	if err != nil {
		return nil, err
	}
	if len(pack.S) > 0 && len(pack.I) > 0 && len(pack.TC) > 0 {
		pack.AC = MsgPackActionCreateOrder
	}
	return pack, nil
}

func GetSettlementId(id, modifier string) string {
	h := md5.New()
	_, _ = io.WriteString(h, id)
	_, _ = io.WriteString(h, modifier)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}

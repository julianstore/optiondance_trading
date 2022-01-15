package core

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/fox-one/mixin-sdk-go"
	"sort"
	"time"
)

type (
	UTXO struct {
		ID              int64               `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		UserID          string              `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		TraceID         string              `gorm:"column:trace_id;type:varchar(45);unique;index:idx_trace_id" json:"trace_id"`
		AssetID         string              `gorm:"column:asset_id;type:varchar(45);index:idx_asset_spent_by" json:"asset_id"`
		Memo            string              `gorm:"column:memo;type:varchar(300)" json:"memo"`
		Amount          float64             `gorm:"column:amount;type:decimal(64,8)" json:"amount"`
		SenderID        string              `gorm:"column:sender_id;type:varchar(45)" json:"sender_id"`
		TransactionHash string              `gorm:"column:transaction_hash;type:varchar(100)" json:"transaction_hash"`
		State           string              `gorm:"column:state;type:varchar(45)" json:"state"`
		SignedBy        string              `gorm:"column:signed_by;type:varchar(100)" json:"signed_by"`
		SignedTx        string              `gorm:"column:signed_tx;type:text" json:"signed_tx"`
		SpentBy         string              `gorm:"column:spent_by;type:varchar(45);index:idx_asset_spent_by" json:"spent_by"`
		CreatedAt       time.Time           `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
		UpdatedAt       time.Time           `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
		Invalid         int8                `gorm:"column:invalid;type:tinyint" json:"invalid"`
		JSONData        string              `gorm:"column:json_data;type:text" json:"json_data"`
		Raw             *mixin.MultisigUTXO `gorm:"-" json:"-,omitempty"`
	}

	UtxoStore interface {
		FindSpentBy(ctx context.Context, assetID, spentBy string) (*UTXO, error)
		Save(ctx context.Context, u *UTXO) (err error)
		List(ctx context.Context, lastId int64, limit int) (utxo []*UTXO, err error)
		ListUnspentV2(ctx context.Context, assetId string, limit int) (utxos []*UTXO, err error)
	}
)

func (UTXO) TableName() string {
	return "utxo"
}

func cmpUTXO(a, b UTXO) int {
	if dur := a.CreatedAt.Sub(b.CreatedAt); dur > 0 {
		return 1
	} else if dur < 0 {
		return -1
	}

	if r := bytes.Compare([]byte(a.TransactionHash[:]), []byte(b.TransactionHash[:])); r != 0 {
		return r
	}
	return -1
}

func SortUTXOs(utxos []UTXO) {
	sort.Slice(utxos, func(i, j int) bool {
		ui, uj := utxos[i], utxos[j]
		return cmpUTXO(ui, uj) < 0
	})
}

func WrapUTXORawData(u *UTXO) (*UTXO, error) {
	var rawUtxo mixin.MultisigUTXO
	err := json.Unmarshal([]byte(u.JSONData), &rawUtxo)
	if err != nil {
		return &UTXO{}, err
	}
	u.Raw = &rawUtxo
	return u, nil
}

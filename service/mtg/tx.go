package mtg

import (
	"context"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/fox-one/mixin-sdk-go"
	"log"
	"option-dance/cmd/config"
	"option-dance/core"
	mixin2 "option-dance/pkg/mixin"
)

const gasAmount = "0.00000001"

func CreateMultiSignTx(ctx context.Context, pack core.MsgPack) (*bot.RawTransaction, error) {

	dapp := config.Cfg.DApp
	memo, err := pack.Pack()
	if err != nil {
		return nil, err
	}
	if len(memo) >= 200 {
		return nil, fmt.Errorf("memo length cannot greater than 200: %v", memo)
	}
	cnbAsset, _ := core.GetAssetIdByCurrency("CNB")
	input := bot.TransferInput{
		AssetId: cnbAsset,
		Amount:  number.FromString(gasAmount),
		TraceId: bot.UuidNewV4().String(),
		Memo:    memo,
		OpponentMultisig: struct {
			Receivers []string
			Threshold int64
		}{
			Receivers: dapp.Receivers,
			Threshold: dapp.Threshold,
		},
	}
	return bot.CreateMultisigTransaction(ctx, &input, dapp.AppID, dapp.SessionID, dapp.PrivateKey, dapp.Pin, dapp.PinToken)
}

func UnlockMutisig(ctx context.Context, raw string) error {

	client, err := mixin2.Client()
	if err != nil {
		log.Panicln(err)
	}
	multisig, err := client.CreateMultisig(context.Background(), mixin.MultisigActionUnlock, raw)
	if err != nil {
		log.Println(err)
	}
	err = client.UnlockMultisig(context.Background(), multisig.RequestID, config.Cfg.DApp.Pin)
	if err != nil {
		log.Println(err)
	}
	return nil
}

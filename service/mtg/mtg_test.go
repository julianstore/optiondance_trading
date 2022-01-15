package mtg

import (
	"context"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
	"log"
	"option-dance/cmd/config"
	"option-dance/core"
	mixin2 "option-dance/pkg/mixin"
	"strconv"
	"testing"
)

const (
	user3 = "5d0b24d1-f2f5-4264-8abc-cdb51b18ee50"
)

func TestMain(m *testing.M) {
	config.InitTestConfig("dev")
	m.Run()
}

func getCreateOrderPaymentUrl(uid, side, price, amount, instrumentName string) (string, error) {

	i, _ := core.ParseInstrument(instrumentName)

	funds := ""
	margin := ""

	if side == "ASK" && i.OptionType == "PUT" {
		funds = decimal.NewFromInt(i.StrikePrice).Mul(decimal.RequireFromString(amount)).String()
		margin = decimal.NewFromInt(i.StrikePrice).Mul(decimal.RequireFromString(amount)).String()
	}

	if side == "BID" && i.OptionType == "PUT" {
		funds = decimal.RequireFromString(price).Mul(decimal.RequireFromString(amount)).String()
	}

	action := core.OrderAction{
		UserId:         uid,
		TraceId:        bot.UuidNewV4().String(),
		QuoteAsset:     "",
		BaseAsset:      "",
		QuoteCurrency:  "pUSD",
		BaseCurrency:   "BTC",
		Side:           side,
		Price:          price,
		Amount:         amount,
		Funds:          funds,
		Type:           "L",
		InstrumentName: instrumentName,
		OptionType:     i.OptionType,
		Margin:         margin,
	}
	request, err := CreateOrderRequest(context.Background(), action)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("mixin://codes/%s", request.Data.CodeId), nil

}

func TestCreateOrderRequest(t *testing.T) {
	for i := 21; i <= 21; i++ {
		instrumentName := fmt.Sprintf("C-pUSD-BTC-1SEP21-%s-P", strconv.Itoa(i))
		side := "ASK"
		price := "1"
		amount := "1"
		url, err := getCreateOrderPaymentUrl(user3, side, price, amount, instrumentName)
		if err != nil {
			t.Error(err)
		}
		t.Log(url)
	}

}

func TestUnlockMutilsig(t *testing.T) {
	client, err := mixin2.Client()
	if err != nil {
		log.Panicln(err)
	}
	raw := "77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d800012a82a2b071b2f98a78f9d26c74fef3c0356ee1aa7eca16d067a05754522438e700010000000000000002000000043b9aca0000019b4664c6282500156195a9f4a09daf113494ec756e64ad6134697ae22db1bd06a6847803ca3d9ab2fb7cde47a0c111f94e383ad82acea6218bc9deff703467da0003fffe0100000000000601be9897f63c000330ea9bab82b534961918b43e7172872cedbef3b8f5899419f97c0b69ec25e5a7862c198e3ef7c5e75fce921432e522e8878767dd88fc5c80221cc02bf5fb61b505a9c6cb2732372a3551c06b9e0206929d4f6522785bc4fde639511d4dd41324890121f2f0ada937183cd999d7ad72884c3d28ddfa198005235e94c747c87ec50003fffe0200000018684b46426f4b46436f4b46506f4b46547055314256454e490000"
	multisig, err := client.CreateMultisig(context.Background(), mixin.MultisigActionUnlock, raw)
	if err != nil {
		log.Println(err)
	}
	err = client.UnlockMultisig(context.Background(), multisig.RequestID, config.Cfg.DApp.Pin)
	println(err)
	if err != nil {
		log.Println(err)
	}
}

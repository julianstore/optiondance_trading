package message

import (
	"option-dance/core"
	"testing"
	"time"
)

func TestMessageBuilder_OrderCreateGroupNotify(t *testing.T) {
	builder := NewMessageBuilder()
	notify, err := builder.OrderCreateGroupNotify(ctx, &core.Order{
		Side:            "ASK",
		OptionType:      "PUT",
		Price:           "12.98",
		RemainingAmount: "1",
		FilledAmount:    "0",
		RemainingFunds:  "1",
		FilledFunds:     "1",
		StrikePrice:     "50000",
		ExpirationDate:  time.Now(),
		QuoteCurrency:   "pUSD",
		BaseCurrency:    "BTC",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(notify)
	messagez := New(client, nil, nil)
	if err = messagez.Send(ctx, notify); err != nil {
		t.Error(err)
	}
	t.Log("success")
}

package market

import (
	"option-dance/cmd/config"
	"option-dance/core"
	"testing"
	"time"
)

var r core.MarketStore

func TestMain(m *testing.M) {
	r = NewRedisStore(config.RedisCli)
	m.Run()
}

func TestRedisStore_Create(t *testing.T) {

	name := "BTC-18JUN21-42000-P"
	bids := `[{
	"side":"ask",
	"price":"123",
	"amount":"12",
	"funds":"12"
}]`
	asks := `[]`

	_, _ = r.UpdateBidAsks(name, bids, asks)
	time.Sleep(2 * time.Second)
	instrumentName, err := r.FindByInstrumentName(name)
	if err != nil {
		t.Error(err)
	}
	t.Log(instrumentName)
	return
}

func TestRedisStore_ListByStrikePrice(t *testing.T) {
	price, err := r.ListByStrikePrice("42000")
	if err != nil {
		t.Error(err)
	}
	t.Log(price)
}

func TestRedisStore_ListMarketStrikePrice(t *testing.T) {
	price, err := r.ListMarketStrikePrice("BID", "CALL")
	if err != nil {
		t.Error(err)
	}
	t.Log(price)
}

package deribit

import (
	"testing"
)

func TestClient_GetInstruments(t *testing.T) {
	instruments, err := C().GetInstruments("BTC", "", false)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(instruments)
}

func TestClient_GetInstrument(t *testing.T) {
	instrument, err := C().GetInstrument("BTC-8OCT21-55000-P")
	if err != nil {
		t.Error(err)
	}
	t.Log(instrument)
}

func TestClient_GetOrderBook(t *testing.T) {
	instrument, err := C().GetOrderBook("BTC-8OCT21-55000-P", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(instrument)
}

package deribit

import (
	"testing"
	"time"
)

func TestClient_GetPositions(t *testing.T) {
	instruments, err := C().GetPositions("BTC", "option")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(instruments)
}

func TestClient_GetPosition(t *testing.T) {
	instruments, err := C().GetPosition("BTC-9JUL21-40000-P")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(instruments)
}

func TestClient_GetTransactionLogResult(t *testing.T) {
	_, err2 := time.Parse("2006-01-02T15:04:05", "2020-01-02T15:04:05")
	if err2 != nil {
		return
	}
	instruments, err := C().GetTransactionLogResult("BTC", "delivery", 0, time.Now().Unix())
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(instruments)
}

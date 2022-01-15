package core

import (
	"strings"
	"testing"
)

func TestGetSettlementId(t *testing.T) {
	id := "2e393a40-6d73-4a7c-8d6f-036774b00ddc"
	i := 1
	for i < 10 {
		i++
		t.Log(GetSettlementId(id, TransferSourcePositionExercised))
		t.Log(GetSettlementId(id, TransferSourceTradeConfirmed))
		t.Log(GetSettlementId(id, TransferSourceOrderCancelled))
		t.Log(GetSettlementId(id, TransferSourcePositionClosedRefund))
		t.Log(GetSettlementId(id, TransferSourcePositionExercisedRefund))
		t.Log("\n")
	}
}

func TestInstrumentExpiryTimestamp(t *testing.T) {
	instrument, _ := ParseInstrument("BTC-7AUG21-40000-P")
	timestamp := InstrumentExpiryTimestamp(instrument.ExpirationDate)
	t.Log(timestamp)
}

func TestSplit(t *testing.T) {
	str := "1213,3123,321312,"
	split1 := strings.Split(str, ",")
	str1 := "123123"
	split2 := strings.Split(str1, ",")
	t.Log(split1)
	t.Log(split2)
	t.Log("ok")
}

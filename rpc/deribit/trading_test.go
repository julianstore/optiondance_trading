package deribit

import (
	"encoding/json"
	"option-dance/pkg/util"
	"testing"
)

func TestClient_GetOpenOrdersByCurrency(t *testing.T) {
	order, err := C().GetOpenOrdersByCurrency("BTC", "", "")
	if err != nil {
		t.Errorf("GetOpenOrdersByCurrency Err %s", err)
	}
	t.Log(order)
}

func TestClient_GetUserTradesByCurrency(t *testing.T) {
	trade, err := C().GetUserTradesByCurrency("BTC", "", "", "", "", 0, false)
	if err != nil {
		t.Errorf("GetOpenOrdersByCurrency Err %s", err)
	}
	t.Log(trade)
}

func TestClient_Buy(t *testing.T) {
	_, err := C().Buy("BTC-29JUN21-30000-P", "0.1", "market")
	if err != nil {
		t.Errorf("GetOpenOrdersByCurrency Err %s", err)
	}
	t.Log()
}

func TestClient_Sell(t *testing.T) {
	buy, err := C().Sell("BTC-9JUL21-42000-P", "0.1", "limit")
	if err != nil {
		t.Errorf("GetOpenOrdersByCurrency Err %s", err)
	}
	marshal, _ := json.Marshal(buy)
	util.JsonPrintS(string(marshal))
	t.Log()
}

func TestClient_Cancel(t *testing.T) {

}

func TestMain(m *testing.M) {
	m.Run()
}

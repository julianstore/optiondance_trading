package position

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/store/deliveryprice"
	"option-dance/store/position"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	join := filepath.Join(wd, "../../config/od_dev_config/node1.yaml")
	config.InitConfig(join, false)
	m.Run()
}

func TestPositionService_Settlement(t *testing.T) {
	SettlementGatherByDate(t, "2021-09-09", core.DeliveryTypePhysical)
	SettlementGatherByDate(t, "2021-09-10", core.DeliveryTypePhysical)
	SettlementGatherByDate(t, "2021-09-13", core.DeliveryTypeCash)
	SettlementGatherByDate(t, "2021-09-13", core.DeliveryTypePhysical)
	t.Log("success")
}

func SettlementGatherByDate(t *testing.T, dateString, deliveryType string) {
	positionTestStore := position.NewPositionTestStore()
	deliveryPriceStore := deliveryprice.New(config.Db)
	service := NewTest(positionTestStore, deliveryPriceStore)
	gather, err := service.SettlementGather(ctx, "", dateString, deliveryType, false, time.Now())
	if err != nil {
		t.Error(err)
	}
	var gotTransfers = gather.Transfers
	//marshal, err := json.Marshal(gotTransfers)
	//if err != nil {
	//	t.Log(err)
	//}
	//err = ioutil.WriteFile(fmt.Sprintf("%s-%s.json", dateString, deliveryType), marshal, 0777)
	//if err != nil {
	//	t.Log(err)
	//}
	var sum1, sum2 decimal.Decimal
	for i := 0; i < len(gotTransfers); i++ {
		t.Logf("userId: %s ,amount: %s ,source: %s", gotTransfers[i].UserId, gotTransfers[i].Amount, gotTransfers[i].Source)
		sum1 = sum1.Add(decimal.RequireFromString(gotTransfers[i].Amount))
	}
	t.Logf("transfer length : %d, total amount :%s", len(gotTransfers), sum1.String())
	transfersData, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s/transfers-%s.json", dateString, strings.ToLower(deliveryType)))
	if err != nil {
		t.Error(err)
	}
	var expectTransfers []*core.Transfer
	if transfersData != nil {
		err = json.Unmarshal(transfersData, &expectTransfers)
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < len(gotTransfers); i++ {
			if !CmpTransfer(gotTransfers[i], expectTransfers[i]) {
				t.Errorf("transfer not matched ,got transfer: %v,expect transfer: %v", gotTransfers[i], expectTransfers[i])
			}
			sum2 = sum2.Add(decimal.RequireFromString(gotTransfers[i].Amount))
		}
		t.Logf("transfer matched ,test %s success,transfer length : %d, total amount :%s", dateString, len(gotTransfers), sum2.String())
	} else {
		t.Error("expect transfer data is nil")
	}
	t.Logf("-----------------------------------%s %s end-----------------------------------------", dateString, deliveryType)
}

func CmpTransfer(t1, t2 *core.Transfer) bool {
	if t1.TransferId == t2.TransferId &&
		t1.UserId == t2.UserId &&
		t1.Opponents == t2.Opponents &&
		t1.Source == t2.Source &&
		t1.Detail == t2.Detail &&
		t1.Amount == t2.Amount {
		return true
	}
	return false
}

func TestToAmountString(t *testing.T) {
	testCase := []struct {
		Input decimal.Decimal
		Want  string
	}{
		{Input: decimal.RequireFromString("0.1"), Want: "0.1"},
		{Input: decimal.RequireFromString("0.12"), Want: "0.12"},
		{Input: decimal.RequireFromString("888888888.1"), Want: "888888888.1"},
		{Input: decimal.RequireFromString("48000.00000000000000321"), Want: "48000"},
		{Input: decimal.RequireFromString("123.98765555555555"), Want: "123.98765555"},
		{Input: decimal.RequireFromString("1.9999999999999"), Want: "1.99999999"},
		{Input: decimal.RequireFromString("1.7654321"), Want: "1.7654321"},
		{Input: decimal.RequireFromString("1.12345678"), Want: "1.12345678"},
	}

	for _, e := range testCase {
		amountString := ToAmountString(e.Input)
		assert.Equal(t, e.Want, amountString)
		t.Logf("want %s, got %s", e.Want, amountString)
	}
}

func TestPositionService_SettlementGather(t *testing.T) {
	positionStore := position.New(config.Db)
	deliveryPriceStore := deliveryprice.New(config.Db)
	service := NewTest(positionStore, deliveryPriceStore)
	dateString := "2021-09-19"
	deliveryType := "CASH"
	gather, err := service.SettlementGather(ctx, "", dateString, deliveryType, false, time.Now())
	if err != nil {
		t.Error(err)
	}
	t.Log(gather)
}

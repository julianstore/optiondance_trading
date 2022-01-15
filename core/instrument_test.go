package core

import (
	"testing"
)

func TestParseInstrument(t *testing.T) {

	testCase := []struct {
		Input string
		Want  OptionInfo
	}{
		{
			Input: "C-pUSD-BTC-24SEP21-80000-P",
			Want: OptionInfo{
				DeliveryType:        DeliveryTypeCash,
				QuoteCurrency:       PUSD,
				BaseCurrency:        BTC,
				ExpirationDate:      InstrumentTime(2021, 9, 24),
				ExpirationTimestamp: InstrumentExpiryTs(2021, 9, 24),
				StrikePrice:         80000,
				OptionType:          OptionTypePUT,
			},
		},
		{
			Input: "P-USDT-BTC-10SEP21-47000-C",
			Want: OptionInfo{
				DeliveryType:        DeliveryTypePhysical,
				QuoteCurrency:       USDT,
				BaseCurrency:        BTC,
				ExpirationDate:      InstrumentTime(2021, 9, 10),
				ExpirationTimestamp: InstrumentExpiryTs(2021, 9, 10),
				StrikePrice:         47000,
				OptionType:          OptionTypeCALL,
			},
		},
	}

	for _, e := range testCase {
		i, err := ParseInstrument(e.Input)
		if err != nil {
			t.Log(err)
		}
		if i.DeliveryType != e.Want.DeliveryType {
			t.Errorf("delivery type error, got %s,want %s", i.DeliveryType, e.Want.DeliveryType)
		}
		if i.QuoteCurrency != e.Want.QuoteCurrency {
			t.Errorf("QuoteCurrency error, got %s,want %s", i.QuoteCurrency, e.Want.QuoteCurrency)
		}
		if i.BaseCurrency != e.Want.BaseCurrency {
			t.Errorf("BaseCurrency error, got %s,want %s", i.BaseCurrency, e.Want.BaseCurrency)
		}
		if i.ExpirationDate.UTC() != e.Want.ExpirationDate.UTC() {
			t.Errorf("ExpirationDate error, got %s,want %s", i.ExpirationDate, e.Want.ExpirationDate)
		}
		if i.ExpirationTimestamp != e.Want.ExpirationTimestamp {
			t.Errorf("ExpirationTimestamp error, got %d,want %d", i.ExpirationTimestamp, e.Want.ExpirationTimestamp)
		}
		if i.StrikePrice != e.Want.StrikePrice {
			t.Errorf("StrikePrice error, got %d,want %d", i.StrikePrice, e.Want.StrikePrice)
		}
		if i.OptionType != e.Want.OptionType {
			t.Errorf("OptionType error, got %s,want %s", i.OptionType, e.Want.OptionType)
		}
	}
}

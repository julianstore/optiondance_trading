package core

import (
	"github.com/stretchr/testify/assert"
	"option-dance/cmd/config"
	"os"
	"path/filepath"
	"testing"
)

func TestGetQuoteBasePair(t *testing.T) {
	var tests = []struct {
		in    string
		quote string
		base  string
	}{
		{in: "BTC", base: currencyMap["BTC"], quote: currencyMap["USDT"]},
		{in: "ETH", base: currencyMap["ETH"], quote: currencyMap["USDT"]},
		{in: "CNB", base: currencyMap["CNB"], quote: currencyMap["USDT"]},
		{in: "XIN", base: currencyMap["XIN"], quote: currencyMap["USDT"]},
		{in: "BTC", base: currencyMap["BTC"], quote: currencyMap["USDT"]},
		{in: "USDT", base: currencyMap["USDT"], quote: currencyMap["USDT"]},
		{in: "SHIB", base: currencyMap["SHIB"], quote: currencyMap["USDT"]},
	}
	for _, e := range tests {
		quoteAsset, baseAsset := GetQuoteBasePair(e.in, e.base)
		assert.Equal(t, e.quote, quoteAsset)
		assert.Equal(t, e.base, baseAsset)
	}
}

func TestGetCurrencyByAsset(t *testing.T) {
	wd, _ := os.Getwd()
	join := filepath.Join(wd, "../config/od_beta_config/test.yaml")
	config.InitConfig(join, false)

	assert.Equal(t, "BTC", GetCurrencyByAsset("dcde18b9-f015-326f-b8b1-5b820a060e44", true))
	assert.Equal(t, "USDT", GetCurrencyByAsset("965e5c6e-434c-3fa9-b780-c50f43cd955c", true))

	assert.Equal(t, SHIB, GetCurrencyByAsset("dcde18b9-f015-326f-b8b1-5b820a060e44", false))
	assert.Equal(t, CNB, GetCurrencyByAsset("965e5c6e-434c-3fa9-b780-c50f43cd955c", false))

	assert.Equal(t, "BTC", GetCurrencyByAsset("c6d0c728-2624-429b-8e0d-d9d19b6592fa", false))
	assert.Equal(t, "USDT", GetCurrencyByAsset("4d8c508b-91c5-375b-92b0-ee702ed2dac5", false))

	assert.Equal(t, "BTC", GetCurrencyByAsset("c6d0c728-2624-429b-8e0d-d9d19b6592fa", true))
	assert.Equal(t, "USDT", GetCurrencyByAsset("4d8c508b-91c5-375b-92b0-ee702ed2dac5", true))
}

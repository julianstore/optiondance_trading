package core

import (
	"fmt"
	"option-dance/cmd/config"
)

const (
	BTC  = "BTC"
	USDT = "USDT"
	XIN  = "XIN"
	ETH  = "ETH"
	CNB  = "CNB"
	SHIB = "SHIB"
	USDC = "USDC"
	PUSD = "pUSD"
)

const GlobalQuoteCurrency = PUSD

var (
	currencyMap = map[string]string{
		"USDT": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
		"BTC":  "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
		"XIN":  "c94ac88f-4671-3976-b60a-09064f1811e8",
		"ETH":  "43d61dcd-e413-450d-80b8-101d5e903357",
		"CNB":  "965e5c6e-434c-3fa9-b780-c50f43cd955c",
		"SHIB": "dcde18b9-f015-326f-b8b1-5b820a060e44",
		"USDC": "9b180ab6-6abe-3dc0-a13f-04169eb34bfa",
		"pUSD": "31d2ea9c-95eb-3355-b65b-ba096853bc18",
	}
)

func GetGlobalQuoteAssetId() string {
	currency, _ := GetAssetIdByCurrency(GlobalQuoteCurrency)
	return currency
}

func GetAssetIdByCurrency(currency string) (string, error) {

	assets := config.Cfg.DApp.Assets

	var assetId string
	if e, found := currencyMap[currency]; found {
		assetId = e
	} else {
		return "", fmt.Errorf("unsupport currency")
	}

	coverAssetId := ""
	switch currency {
	case USDT:
		coverAssetId = assets.USDT
		break
	case USDC:
		coverAssetId = assets.USDC
		break
	case PUSD:
		coverAssetId = assets.PUSD
		break
	case BTC:
		coverAssetId = assets.BTC
		break
	case CNB:
		coverAssetId = assets.CNB
		break
	case XIN:
		coverAssetId = assets.XIN
		break
	case ETH:
		coverAssetId = assets.ETH
		break
	default:
		return "", fmt.Errorf("unsupport currency")
	}
	//assetId covered
	if len(coverAssetId) > 0 {
		assetId = coverAssetId
	}
	return assetId, nil
}

// GetCurrencyByAsset covered means if use the covered asset configured in config file dapp.assets
func GetCurrencyByAsset(assetId string, covered bool) (currency string) {
	assets := config.Cfg.DApp.Assets
	switch assetId {
	case currencyMap[USDT]:
		currency = USDT
		break
	case currencyMap[USDC]:
		currency = USDC
		break
	case currencyMap[BTC]:
		currency = BTC
		break
	case currencyMap[PUSD]:
		currency = PUSD
		break
	case currencyMap[SHIB]:
		if covered && assets.BTC == assetId {
			currency = BTC
		} else {
			currency = SHIB
		}
		break
	case currencyMap[CNB]:
		if covered && assets.PUSD == assetId {
			currency = PUSD
		} else if covered && assets.USDC == assetId {
			currency = USDC
		} else if covered && assets.USDT == assetId {
			currency = USDT
		} else {
			currency = CNB
		}
	case currencyMap[XIN]:
		currency = XIN
		break
	default:
		currency = ""
		break
	}
	return
}

func GetQuoteBasePair(quoteCurrency, baseCurrency string) (quoteAsset, baseAsset string) {
	quoteAsset, _ = GetAssetIdByCurrency(quoteCurrency)
	baseAsset, _ = GetAssetIdByCurrency(baseCurrency)
	return
}

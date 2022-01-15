package deribit

import (
	"encoding/json"
)

const (
	urlGetInstruments = "/public/get_instruments"
	urlGetInstrument  = "/public/get_instrument"
	urlGetOrderBook   = "/public/get_order_book"
	urlGetIndexPrice  = "/public/get_index_price"
)

//https://docs.deribit.com/?shell#public-get_instrument
func (c *Client) GetInstruments(currency, kind string, expired bool) (o InstrumentListResult, err error) {
	res, err := c.Get(urlGetInstruments, map[string]interface{}{
		"currency": currency,
		"kind":     kind,
		"expired":  expired,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

func (c *Client) GetInstrument(instrumentName string) (o InstrumentResult, err error) {
	res, err := c.Get(urlGetInstrument, map[string]interface{}{
		"instrument_name": instrumentName,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

func (c *Client) GetOrderBook(instrumentName string, depth int) (o OrderBookResult, err error) {
	res, err := c.Get(urlGetOrderBook, map[string]interface{}{
		"instrument_name": instrumentName,
		"depth":           depth,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

func (c *Client) GetIndexPrice(indexName string) (o IndexPriceResult, err error) {
	res, err := c.Get(urlGetIndexPrice, map[string]interface{}{
		"index_name": indexName,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

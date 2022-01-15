package deribit

import (
	"encoding/json"
	"strconv"
)

const (
	urlGetPositions      = "/private/get_positions"
	urlGetPosition       = "/private/get_position"
	urlGetTransactionLog = "/private/get_transaction_log"
)

//https://docs.deribit.com/#private-get_positions
func (c *Client) GetPositions(currency, kind string) (o PositionListResult, err error) {
	res, err := c.Get(urlGetPositions, map[string]interface{}{
		"currency": currency,
		"kind":     kind,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

func (c *Client) GetPosition(instrumentName string) (o GetPositionResult, err error) {
	res, err := c.Get(urlGetPosition, map[string]interface{}{
		"instrument_name": instrumentName,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

func (c *Client) GetTransactionLogResult(currency, query string, startTs, endTs int64) (o GetTransactionLogResult, err error) {
	res, err := c.Get(urlGetTransactionLog, map[string]interface{}{
		"currency":        currency,
		"start_timestamp": strconv.Itoa(int(startTs)),
		"end_timestamp":   strconv.Itoa(int(endTs)),
		"query":           query,
	})
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(res, &o)
	return o, err
}

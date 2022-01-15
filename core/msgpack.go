package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/MixinNetwork/go-number"
	"github.com/asaskevich/govalidator"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/ugorji/go/codec"
)

const (
	//Create Order
	MsgPackActionCreateOrder = "CREATE_ORDER"
	//Cancel order
	MsgPackActionCancelOrder = "CANCEL_ORDER"
	//Exercise
	MsgPackActionExercise = "EXERCISE"
	//Automatic exercise
	MsgPackActionAutoExercise = "AUTO_EXERCISE"
	//Settlement and delivery
	MsgPackActionSettlement = "SETTLEMENT"
	//Automatic cancellation on expiry date
	MsgPackActionExpiringCancelOrder = "EXPIRING_CANCEL_ORDER"
	//Notice before expiry
	MsgPackActionExpiringNotify = "EXPIRING_NOTIFY"
	//sync deribit delivery price
	MsgPackSyncDbDeliveryPrice = "SYNC_DB_DELIVERY_PRICE"

	PackOrderTypeLimit  = "L"
	PackOrderTypeMarket = "M"
)

const (
	OptionAmountPlaces    = 1
	PutOptionPricePlaces  = 2
	CallOptionPricePlaces = 4

	place1 = `^(([1-9]{1}\d*)|(0{1}))(\.\d{0,1})?$`
	place2 = `^(([1-9]{1}\d*)|(0{1}))(\.\d{0,2})?$`
	place3 = `^(([1-9]{1}\d*)|(0{1}))(\.\d{0,3})?$`
	place4 = `^(([1-9]{1}\d*)|(0{1}))(\.\d{0,4})?$`
	place5 = `^(([1-9]{1}\d*)|(0{1}))(\.\d{0,5})?$`
)

type MsgPack struct {
	U  string // user
	S  string // side
	T  string // type
	I  string // instrument
	P  string // price
	O  string // order
	A  string // amount
	TC string // trace id
	Q  string // quote
	QC string // quote currency
	B  string // base
	BC string // base currency
	M  string // margin
	AC string // action
}

func CheckDecimalPlaces(f string, places int) (passed bool, err error) {
	// Determine the number of decimal points after the decimal point
	_, err = decimal.NewFromString(f)
	if err != nil {
		return false, err
	}
	switch places {
	case 1:
		return regexp.MatchString(place1, f)
	case 2:
		return regexp.MatchString(place2, f)
	case 3:
		return regexp.MatchString(place3, f)
	case 4:
		return regexp.MatchString(place4, f)
	case 5:
		return regexp.MatchString(place5, f)
	default:
		return false, fmt.Errorf("not support places %d", places)
	}
}

func (p MsgPack) Pack() (memo string, err error) {
	b := make([]byte, 140)
	handle := new(codec.MsgpackHandle)
	encoder := codec.NewEncoderBytes(&b, handle)
	err = encoder.Encode(p)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (p MsgPack) CheckCreateOrderPack() (err error) {
	err = errors.New("order msg pack format error")
	//check order type
	if !govalidator.IsIn(p.T, OrderTypeLimit, OrderTypeMarket, PackOrderTypeLimit, PackOrderTypeMarket) {
		return err
	}
	//check order id
	if p.TC == "" {
		return err
	}
	_, err = uuid.FromString(p.TC)
	if err != nil {
		return err
	}
	//check expiry date
	i, err := ParseInstrument(p.I)
	if err != nil {
		return fmt.Errorf("utxo:create_order:instrument format error")
	}
	price := number.FromString(p.P)
	amount := number.FromString(p.A)

	var expectPricePlaces = PutOptionPricePlaces
	if i.OptionType == OptionTypeCALL {
		expectPricePlaces = CallOptionPricePlaces
	}

	//check price, 0~99999
	if !(price.Cmp(number.Zero()) > 0) || price.Cmp(number.FromString("99999")) > 0 {
		return fmt.Errorf("price incorrect: %v", p.P)
	}
	//check price places
	pricePlacesPassed, err := CheckDecimalPlaces(p.P, expectPricePlaces)
	if err != nil {
		return err
	}
	if !pricePlacesPassed {
		return fmt.Errorf("price format incorrect: %v, expect %d places", p.P, expectPricePlaces)
	}

	//check amount, 0~99999
	if !(amount.Cmp(number.Zero()) > 0) || amount.Cmp(number.FromString("99999")) > 0 {
		return fmt.Errorf("amount incorrect: %v", p.P)
	}

	//check amount places
	amountPlacesPassed, err := CheckDecimalPlaces(p.A, OptionAmountPlaces)
	if err != nil {
		return err
	}
	if !amountPlacesPassed {
		return fmt.Errorf("amount format incorrect: %v, expect %d places", p.P, OptionAmountPlaces)
	}

	//check margin in sell put/sell call
	margin := decimal.RequireFromString(p.M)
	var expectMargin decimal.Decimal
	if p.S == "A" || p.S == PageSideAsk {
		if i.OptionType == OptionTypePUT {
			expectMargin = decimal.NewFromInt(i.StrikePrice).Mul(decimal.RequireFromString(p.A))
		} else if i.OptionType == OptionTypeCALL {
			expectMargin = decimal.RequireFromString(p.A)
		} else {
			return fmt.Errorf("wrong option Type %s", i.OptionType)
		}
		if !margin.Equal(expectMargin) {
			return fmt.Errorf("wrong order margin: %v", margin)
		}
	}
	return nil
}

func (p MsgPack) CheckCreateOrderPackWithUtxo(u *UTXO) (err error) {
	if err = p.CheckCreateOrderPack(); err != nil {
		return err
	}
	//check funds and utxo.amount
	if p.S == "B" || p.S == PageSideBid {
		funds := decimal.RequireFromString(p.P).Mul(decimal.RequireFromString(p.A))
		if !funds.Equal(decimal.NewFromFloat(u.Amount)) {
			return fmt.Errorf("create order funds not correct ,price %s,amount %s,funds %s", p.P, p.A, funds)
		}
	} else if p.S == "A" || p.S == PageSideAsk {
		//check margin and utxo.amount
		i, _ := ParseInstrument(p.I)
		var expectMargin decimal.Decimal
		if i.OptionType == OptionTypePUT {
			expectMargin = decimal.NewFromInt(i.StrikePrice).Mul(decimal.RequireFromString(p.A))
		} else {
			expectMargin = decimal.RequireFromString(p.A)
		}
		if !expectMargin.Equal(decimal.NewFromFloat(u.Amount)) {
			return fmt.Errorf("utxo wrong order margin: %v", u)
		}
	} else {
		return fmt.Errorf("utxo wrong order option side: %v", u)
	}
	return nil
}

type OrderAction struct {
	UserId         string `json:"user_id"`
	TraceId        string `json:"trace_id"`
	QuoteAsset     string `json:"quote_asset"`
	BaseAsset      string `json:"base_asset"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseCurrency   string `json:"base_currency"`
	Side           string `json:"side"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Funds          string `json:"funds"`
	Type           string `json:"type"`
	InstrumentName string `json:"instrument_name"`
	OptionType     string `json:"option_type"`
	Margin         string `json:"margin"`
}

func DecryptAction(data string) (*MsgPack, error) {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		payload, err = base64.URLEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
	}
	var action MsgPack
	cdc := new(codec.MsgpackHandle)
	decoder := codec.NewDecoderBytes(payload, cdc)
	err = decoder.Decode(&action)
	if err != nil {
		return nil, err
	}
	switch action.T {
	case "L":
		action.T = OrderTypeLimit
	case "M":
		action.T = OrderTypeMarket
	}
	switch action.S {
	case "A":
		action.S = PageSideAsk
	case "B":
		action.S = PageSideBid
	}
	return &action, nil
}

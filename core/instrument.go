package core

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	WrongInstrumentFormatErr = errors.New("wrong instrument format")
)

const (
	OptionTypePUT        = "PUT"
	OptionTypeCALL       = "CALL"
	DeliveryTypeCash     = "CASH"
	DeliveryTypePhysical = "PHYSICAL"

	// 8：00 ~ 16：00
	exercisePeriod = 8 * 60 * 60
	//16:00 ~ 16:30
	expiryToSettlementPeriod = 30 * 60
)

const (
	idxDeliveryType = iota + 1
	idxQuoteCurrency
	idxBaseCurrency
	idxDay
	idxMonth
	idxYear
	idxStrikePrice
	idxOptionType
)

var (
	monthEnMap = map[string]string{
		"01": "JAN", "02": "FEB", "03": "MAR", "04": "APR", "05": "MAY", "06": "JUN",
		"07": "JUL", "08": "AUG", "09": "SEP", "10": "OCT", "11": "NOV", "12": "DEC",
	}
)

// To determine whether it can be exercised
func Exercisable(exercisedTs, expiryTs int64) bool {
	settlementTs := expiryTs + expiryToSettlementPeriod
	return exercisedTs < settlementTs && exercisedTs > expiryTs-exercisePeriod
}

func ExercisableContract(exercisedTs, expiryTs int64) bool {
	return exercisedTs > expiryTs-exercisePeriod
}

// To determine whether an instrument is expired
func Expired(createAtTs, expiryTs int64) bool {
	return createAtTs >= expiryTs
}

func ExpiredPosition(createAtTs, expiryTs int64) bool {
	return createAtTs >= expiryTs+expiryToSettlementPeriod
}

type OptionInfo struct {
	DeliveryType        string
	QuoteCurrency       string
	BaseCurrency        string
	ExpirationDate      time.Time
	ExpirationTimestamp int64
	StrikePrice         int64
	OptionType          string
}

// ParseInstrument eg. format C-pUSD-BTC-24SEP21-80000-P
func ParseInstrument(i string) (info OptionInfo, err error) {
	monthMap := map[string]string{"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04", "MAY": "05", "JUN": "06", "JUL": "07", "AUG": "08", "SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12"}
	reg := `^(C|P)-(pUSD|USDT|BTC)-(BTC|XIN|ETH)-([\d]{1,2})(JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)([\d]{2})-([\d]+)-(C|P)$`
	instrumentReg, err := regexp.Compile(reg)
	if err != nil {
		return info, err
	}
	params := instrumentReg.FindStringSubmatch(i)

	var deliveryType string
	if params[idxDeliveryType] == "P" {
		deliveryType = DeliveryTypePhysical
	} else if params[idxDeliveryType] == "C" {
		deliveryType = DeliveryTypeCash
	} else {
		return info, fmt.Errorf("instrument deliveryType error: %s", i)
	}

	var day = params[idxDay]
	if len(day) == 1 {
		day = "0" + day
	}
	timeString := fmt.Sprintf("20%s-%s-%sT16:00:00+08:00", params[idxYear], monthMap[params[idxMonth]], day)
	expirationDate, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return info, err
	}
	strikePrice, err := strconv.Atoi(params[idxStrikePrice])
	if err != nil {
		return info, err
	}
	var optionType string
	if params[idxOptionType] == "P" {
		optionType = OptionTypePUT
	} else if params[idxOptionType] == "C" {
		optionType = OptionTypeCALL
	} else {
		return info, fmt.Errorf("instrument optionType error: %s", i)
	}
	if len(params) == 9 {
		info := OptionInfo{
			DeliveryType:        deliveryType,
			QuoteCurrency:       params[idxQuoteCurrency],
			BaseCurrency:        params[idxBaseCurrency],
			ExpirationDate:      expirationDate,
			ExpirationTimestamp: expirationDate.UTC().Unix(),
			StrikePrice:         int64(strikePrice),
			OptionType:          optionType,
		}
		return info, nil
	} else {
		return info, fmt.Errorf("instrument format error: %s", i)
	}
}

func InstrumentDate(time time.Time) string {
	format := time.Format("06-01-02")
	split := strings.Split(format, "-")
	monthMap := map[string]string{"01": "JAN", "02": "FEB", "03": "MAR", "04": "APR", "05": "MAY", "06": "JUN", "07": "JUL", "08": "AUG", "09": "SEP", "10": "OCT", "11": "NOV", "12": "DEC"}
	day := split[2]
	if strings.HasPrefix(day, "0") {
		day = day[1:]
	}
	return fmt.Sprintf("%s%s%s", day, monthMap[split[1]], split[0])
}

func InstrumentExpiryTimestamp(t time.Time) int64 {
	date := time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, time.UTC)
	return date.Unix()
}

func InstrumentExpiryTs(year, month, day int) int64 {
	return InstrumentTime(year, month, day).Unix()
}

func InstrumentTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 8, 0, 0, 0, time.UTC)
}

func ParseOptionType(instrument string) string {
	var optionType = OptionTypePUT
	if strings.HasSuffix(instrument, "-C") {
		optionType = OptionTypeCALL
	}
	return optionType
}

// ToOdInstrument eg. format C-pUSD-BTC-24SEP21-80000-P
func ToOdInstrument(deliveryType, quoteCurrency, baseCurrency, strikePrice, optionType string, expiryDate time.Time) string {
	var (
		dt = "C"
		qc = "pUSD"
		bc = "BTC"
		ep = ""
		sp = strikePrice
		ot = "P"
	)
	if deliveryType == DeliveryTypePhysical {
		dt = "P"
	}
	if quoteCurrency != "" {
		qc = quoteCurrency
	}
	if baseCurrency != "" {
		bc = baseCurrency
	}
	if optionType == OptionTypeCALL {
		ot = "C"
	}
	//expiryDate
	format := expiryDate.Format("020106")
	day := format[:2]
	if strings.HasPrefix(day, "0") {
		day = day[1:]
	}
	ep = fmt.Sprintf("%s%s%s", day, monthEnMap[format[2:4]], format[4:])
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s", dt, qc, bc, ep, sp, ot)
}

// ToDbInstrument eg. format BTC-24SEP21-80000-P
func ToDbInstrument(baseCurrency, strikePrice, optionType string, expiryDate time.Time) string {
	var (
		bc = "BTC"
		ep = ""
		sp = strikePrice
		ot = "P"
	)
	if baseCurrency != "" {
		bc = baseCurrency
	}
	if optionType == OptionTypeCALL {
		ot = "C"
	}
	//expiryDate
	format := expiryDate.Format("020106")
	day := format[:2]
	if strings.HasPrefix(day, "0") {
		day = day[1:]
	}
	ep = fmt.Sprintf("%s%s%s", day, monthEnMap[format[2:4]], format[4:])
	return fmt.Sprintf("%s-%s-%s-%s", bc, ep, sp, ot)
}

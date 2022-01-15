package statistics

import (
	"context"
	"github.com/shopspring/decimal"
	"option-dance/core"
	"option-dance/store"
	"strconv"
)

func New(positionStore core.PositionStore, tradeStore core.TradeStore) core.StatisticsService {
	return &service{positionStore: positionStore, tradeStore: tradeStore}
}

type service struct {
	positionStore core.PositionStore
	tradeStore    core.TradeStore
}

func (s *service) ListMonthlyPremium(ctx context.Context, uid, side string, year int) (annualProfit core.AnnualProfit, err error) {
	monthlyProfitList, err := s.tradeStore.ListMonthlyPremium(ctx, uid, side, year)
	if store.CheckErr(err) != nil {
		return annualProfit, err
	}
	asset := core.GlobalQuoteCurrency
	annualProfit.MonthlyProfitList, annualProfit.TotalAmount = s.wrapper12Month(monthlyProfitList, year, asset)
	annualProfit.Year = year
	annualProfit.Asset = asset
	return
}

func (s *service) ListMonthlyUnderlying(ctx context.Context, uid string, year int) (annualProfit core.AnnualProfit, err error) {
	monthlyProfitList, err := s.positionStore.ListMonthlyUnderlying(ctx, uid, year)
	if store.CheckErr(err) != nil {
		return annualProfit, err
	}
	asset := "BTC"
	annualProfit.MonthlyProfitList, annualProfit.TotalAmount = s.wrapper12Month(monthlyProfitList, year, asset)
	annualProfit.Year = year
	annualProfit.Asset = asset
	return
}

func (s *service) wrapper12Month(monthlyProfitList []core.MonthlyProfit, year int, asset string) (result []core.MonthlyProfit, totalAmount string) {
	var total decimal.Decimal
	result = make([]core.MonthlyProfit, 12)
	for _, e := range monthlyProfitList {
		monthString := e.Month[len(e.Month)-2:]
		if monthString[:1] == "0" {
			monthString = monthString[1:]
		}
		month, _ := strconv.Atoi(monthString)
		e.Month = strconv.Itoa(month)
		e.Asset = asset
		e.Year = year
		e.Amount = decimal.RequireFromString(e.Amount).StringFixed(4)
		if decimal.RequireFromString(e.Amount).Cmp(decimal.Zero) == 0 {
			e.Amount = "0"
		}
		result[month-1] = e
		total = total.Add(decimal.RequireFromString(e.Amount))
	}
	totalAmount = total.String()
	for i, e := range result {
		if e.Month == "" {
			result[i] = core.MonthlyProfit{
				Year:   year,
				Month:  strconv.Itoa(i + 1),
				Amount: "0",
				Asset:  asset,
			}
		}
	}
	return
}

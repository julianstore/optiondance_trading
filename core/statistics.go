package core

import "context"

type (
	MonthlyProfit struct {
		Year   int    `json:"year"`
		Month  string `json:"month"`
		Amount string `json:"amount"`
		Asset  string `json:"asset"`
	}

	AnnualProfit struct {
		Asset             string          `json:"asset"`
		TotalAmount       string          `json:"total_amount"`
		Year              int             `json:"year"`
		MonthlyProfitList []MonthlyProfit `json:"monthly_profit_list"`
	}

	StatisticsService interface {
		ListMonthlyPremium(ctx context.Context, uid, side string, year int) (annualProfit AnnualProfit, err error)
		ListMonthlyUnderlying(ctx context.Context, uid string, year int) (annualProfit AnnualProfit, err error)
	}
)

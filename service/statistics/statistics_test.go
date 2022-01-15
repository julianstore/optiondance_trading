package statistics

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/store"
	"testing"
)

var (
	db                *gorm.DB
	ctx               context.Context
	statisticsService core.StatisticsService
	positionStore     core.PositionStore
)

func TestService_ListMonthlyPremium(t *testing.T) {
	underlying, err := statisticsService.ListMonthlyPremium(ctx, "019308a5-e1a9-427c-af2a-e05093beedaa", "ASK", 2021)
	if store.CheckErr(err) != nil {
		t.Error(err)
	}
	t.Log(underlying)
}

func TestService_ListMonthlyUnderlying(t *testing.T) {
	underlying, err := statisticsService.ListMonthlyUnderlying(ctx, "019308a5-e1a9-427c-af2a-e05093beedaa", 2021)
	if store.CheckErr(err) != nil {
		t.Error(err)
	}
	t.Log(underlying)
}

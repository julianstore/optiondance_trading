package trade

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
)

func New(db *gorm.DB) core.TradeStore {
	return &tradeStore{db: db}
}

type tradeStore struct {
	db *gorm.DB
}

func (s *tradeStore) ListMonthlyPremium(ctx context.Context, uid, side string, year int) (monthlyProfit []core.MonthlyProfit, err error) {
	querySql := `SELECT date_format(t.created_at,'%Y%m')as month,
					sum(t.price * t.amount) as amount
					FROM trade t
					where t.user_id = ?
					and t.side = ? and year(created_at) = ?
					group by month`
	if err = s.db.Model(core.Position{}).
		Raw(querySql, uid, side, year).Find(&monthlyProfit).Error; err != nil {
		return nil, err
	}
	return
}

func (s *tradeStore) ListOrderTrades(userId string, orderId string, side string) (tradeList []*core.Trade, err error) {
	var queryString string
	if side == "BID" {
		queryString = "user_id = ? AND bid_order_id = ? AND side = 'BID'"
	} else {
		queryString = "user_id = ? AND ask_order_id = ? AND side = 'ASK'"
	}
	if err := s.db.Model(core.Trade{}).
		Where(queryString, userId, orderId).
		Order("id").Find(&tradeList).Error; err != nil {
		return nil, err
	}
	return tradeList, nil
}

package order

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/pkg/util"
	"option-dance/store"
	"strconv"
	"time"
)

func New(db *gorm.DB) core.OrderStore {
	return &orderStore{db: db}
}

type orderStore struct {
	db *gorm.DB
}

func (s *orderStore) ListInstrumentNamesByExpiryDate(ctx context.Context, time time.Time) (names []string, err error) {
	date := core.InstrumentDate(time)
	if len(date) == 6 {
		date = "-" + date
	}
	if err := s.db.Model(core.Order{}).WithContext(ctx).
		Raw("SELECT instrument_name FROM `order` where instrument_name like ? AND order_status = ? group by instrument_name", "%"+date+"%", core.OrderStatusOpen).
		Scan(&names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

func (s *orderStore) ListInstrumentNamesBeforeExpiryDate(ctx context.Context, time time.Time) (names []string, err error) {
	ts := core.InstrumentExpiryTimestamp(time)
	if err := s.db.Model(core.Order{}).WithContext(ctx).
		Raw("SELECT instrument_name FROM `order` where expiration_timestamp <= ? AND order_status = ? group by instrument_name", ts, core.OrderStatusOpen).
		Scan(&names).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return names, nil
}

func (s *orderStore) ListOrdersByINameAndStatus(ctx context.Context, instrumentName string, status int64) (orders []*core.Order, err error) {

	if err := s.db.Model(core.Order{}).WithContext(ctx).
		Where("instrument_name = ? AND order_status = ?", instrumentName, status).Find(&orders).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) ListOrdersBeforeExpiryDate(ctx context.Context, status int64, time time.Time) (orders []*core.Order, err error) {

	expiryTs := core.InstrumentExpiryTimestamp(time)
	if err := s.db.Model(core.Order{}).WithContext(ctx).
		Where("  expiration_timestamp <= ? AND order_status = ?  ", expiryTs, status).Find(&orders).Error; store.CheckErr(err) != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderStore) ListOrdersBeforeTime(ctx context.Context, status int64, side string, time time.Time) (orders []*core.Order, err error) {
	sideSet := make([]string, 0)
	if side == "ALL" {
		sideSet = append(sideSet, core.PageSideAsk, core.PageSideBid)
	} else {
		sideSet = append(sideSet, side)
	}
	if err := s.db.Model(core.Order{}).WithContext(ctx).
		Where("created_at <= ? AND side in (?) AND order_status = ? ", time, sideSet, status).
		Find(&orders).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) FindByOrderId(ctx context.Context, orderId string) (order *core.Order, err error) {
	if err = s.db.Model(core.Order{}).WithContext(ctx).
		Where("order_id=?", orderId).First(&order).Error; store.CheckErr(err) != nil {
		return order, err
	}
	return order, nil
}

func (s *orderStore) ListOrderByIds(ctx context.Context, orderIds []string) (orders []core.Order, err error) {
	if err = s.db.Model(core.Order{}).Where("order_id in ?", orderIds).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) ListOrdersByPage(current int64, size int64, uid, status, order string) (list []*core.Order, total int64, pages int64) {
	if status == "open" {
		s.db.Model(core.Order{}).Where("user_id =? and order_status = 10", uid).Count(&total)
		pages = util.GetPages(total, size)
		s.db.Model(core.Order{}).Where("user_id=?  and order_status = 10", uid).
			Offset(int((current - 1) * size)).Limit(int(size)).Order("created_at desc").Find(&list)
	} else {
		s.db.Model(core.Order{}).Where("user_id =? and order_status > 20", uid).Count(&total)
		pages = util.GetPages(total, size)
		s.db.Model(core.Order{}).Where("user_id=?  and order_status > 20", uid).
			Offset(int((current - 1) * size)).Limit(int(size)).Order("created_at desc").Find(&list)
	}
	for _, e := range list {
		e.OrderNumString = strconv.Itoa(int(e.OrderNum))
	}
	return
}

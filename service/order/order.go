package order

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/match/engine"
	"option-dance/pkg/util"
	"option-dance/store"
	"strconv"
	"time"
)

const (
	OrderActionTypeUser = "USER"
)

func New(orderStore core.OrderStore,
	tradeStore core.TradeStore,
	actionStore core.ActionStore,
) core.OrderService {
	return &orderService{orderStore: orderStore, tradeStore: tradeStore, actionStore: actionStore}
}

type orderService struct {
	actionStore core.ActionStore
	orderStore  core.OrderStore
	tradeStore  core.TradeStore
}

func (o *orderService) ListOpenOrdersByExpiryDate(ctx context.Context, time time.Time) (orders []*core.Order, err error) {

	names, err := o.orderStore.ListInstrumentNamesByExpiryDate(ctx, time)
	if err != nil {
		return nil, err
	}

	for _, e := range names {
		openOrders, err := o.orderStore.ListOrdersByINameAndStatus(ctx, e, core.OrderStatusOpen)
		if err != nil {
			return nil, err
		}
		orders = append(orders, openOrders...)
	}
	return orders, nil
}

func (o *orderService) ListOpenOrdersBeforeExpiryDate(ctx context.Context, time time.Time) (openOrders []*core.Order, err error) {

	openOrders, err = o.orderStore.ListOrdersBeforeExpiryDate(ctx, core.OrderStatusOpen, time)
	if err != nil {
		return nil, err
	}

	return openOrders, nil
}

func (o *orderService) ListOpenOrdersBeforeTime(ctx context.Context, side string, time time.Time) (openOrders []*core.Order, err error) {

	openOrders, err = o.orderStore.ListOrdersBeforeTime(ctx, core.OrderStatusOpen, side, time)
	if err != nil {
		return nil, err
	}
	return openOrders, nil
}

func (o *orderService) CreateOrderAction(ctx context.Context, actionOrder *core.Order, userId, actionType string, createdAt time.Time) error {
	if !number.FromString(actionOrder.FilledFunds).IsZero() || !number.FromString(actionOrder.FilledAmount).IsZero() {
		zap.L().Panic("CreateOrderAction panic", zap.String("userId", userId), zap.Any("order", actionOrder))
	}
	orderNum, err := util.GenSnowflakeIdInt64()
	if err != nil {
		return err
	}
	option, err := core.ParseInstrument(actionOrder.InstrumentName)
	if err != nil {
		return err
	}
	order := core.Order{
		OrderID:             actionOrder.OrderID,
		OrderNum:            orderNum,
		OrderType:           actionOrder.OrderType,
		Side:                actionOrder.Side,
		QuoteAssetID:        actionOrder.QuoteAssetID,
		BaseAssetID:         actionOrder.BaseAssetID,
		QuoteCurrency:       actionOrder.QuoteCurrency,
		BaseCurrency:        actionOrder.BaseCurrency,
		InstrumentName:      actionOrder.InstrumentName,
		StrikePrice:         strconv.Itoa(int(option.StrikePrice)),
		ExpirationDate:      option.ExpirationDate,
		ExpirationTimestamp: option.ExpirationTimestamp,
		OptionType:          option.OptionType,
		DeliveryType:        option.DeliveryType,
		Price:               actionOrder.Price,
		RemainingAmount:     actionOrder.RemainingAmount,
		FilledAmount:        actionOrder.FilledAmount,
		RemainingFunds:      actionOrder.RemainingFunds,
		Margin:              actionOrder.Margin,
		FilledFunds:         actionOrder.FilledFunds,
		CreatedAt:           createdAt,
		OrderStatus:         core.OrderStatusOpen,
		UserID:              userId,
	}
	action := core.Action{
		OrderId:   order.OrderID,
		Action:    engine.OrderActionCreate,
		CreatedAt: createdAt,
	}
	err = config.Db.Transaction(func(tx *gorm.DB) error {
		_, err = checkOrderState(ctx, tx, order.OrderID)
		if err != nil {
			return err
		}
		var count int64
		if err := tx.Model(core.Order{}).Where("order_id = ?", order.OrderID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err = tx.Model(core.Order{}).Create(&order).Error; err != nil {
				return err
			}
			if err = tx.Model(core.Action{}).Create(&action).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (o *orderService) CancelOrderAction(ctx context.Context, orderId string, createdAt time.Time) error {
	action := core.Action{
		OrderId:   orderId,
		Action:    engine.OrderActionCancel,
		CreatedAt: createdAt,
	}
	err := config.Db.Transaction(func(tx *gorm.DB) error {
		exist, err := checkActionExistence(ctx, tx, orderId, action.Action)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
		o, err := checkOrderState(ctx, tx, action.OrderId)
		if err != nil {
			return err
		}
		if o.ID == 0 || o.OrderStatus == core.OrderStatusCancelled {
			return nil
		}
		if err = tx.Model(core.Action{}).Create(&action).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func checkOrderState(ctx context.Context, tx *gorm.DB, orderId string) (order core.Order, err error) {
	if err = tx.Where("order_id=?", orderId).WithContext(ctx).First(&order).Error; store.CheckErr(err) != nil {
		return order, err
	}
	return order, nil
}

func (o *orderService) ListPendingActions(ctx context.Context, checkpoint int, limit int) (actions []*core.Action, err error) {

	actions, err = o.actionStore.ListOrderActions(ctx, checkpoint, limit)
	if err != nil {
		return nil, err
	}
	var (
		orderFilters = make(map[string]bool)
		orderIds     = make([]string, 0)
		orders       = make([]core.Order, 0)
		orderMap     = make(map[string]core.Order)
	)

	for _, e := range actions {
		filterKey := e.GetFilterKey()
		if orderFilters[filterKey] {
			continue
		}
		orderFilters[filterKey] = true
		orderIds = append(orderIds, e.OrderId)
	}

	orders, err = o.orderStore.ListOrderByIds(ctx, orderIds)
	if err != nil {
		return nil, err
	}

	for _, e := range orders {
		orderMap[e.OrderID] = e
	}

	for _, a := range actions {
		a.Order = orderMap[a.OrderId]
	}
	return actions, nil
}

func checkActionExistence(ctx context.Context, tx *gorm.DB, orderId, action string) (exist bool, err error) {
	var count int64
	err = tx.WithContext(ctx).Model(core.Action{}).Where("order_id=? AND action=?", orderId, action).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (o *orderService) Page(current int64, size int64, uid string, status, order string) (list []*core.Order, total int64, pages int64) {
	return o.orderStore.ListOrdersByPage(current, size, uid, status, order)
}

func (o *orderService) FindByOrderId(ctx context.Context, userId, orderId string) (orders *core.Order, err error) {

	order, err := o.orderStore.FindByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}
	trades, err := o.tradeStore.ListOrderTrades(userId, orderId, order.Side)
	if err != nil {
		return nil, err
	}
	order.TradeList = trades
	order.OrderNumString = strconv.Itoa(int(order.OrderNum))
	return order, nil
}

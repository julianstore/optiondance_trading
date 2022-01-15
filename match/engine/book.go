package engine

import (
	"context"
	"encoding/json"
	"github.com/MixinNetwork/go-number"
	"log"
	"option-dance/core"
	"option-dance/match/cache"
	"option-dance/pkg/util"
	"time"
)

const (
	OrderActionCreate = "CREATE"
	OrderActionCancel = "CANCEL"

	EventQueueSize = 8192
)

type TransactCallback func(taker, maker *core.EngineOrder, amount number.Integer) string
type CancelCallback func(order *core.EngineOrder, cancelledAt time.Time)

type OrderEvent struct {
	Order  *core.EngineOrder
	Action *core.Action
}

type Book struct {
	market        string
	createIndex   map[string]bool
	cancelIndex   map[string]bool
	transact      TransactCallback
	cancel        CancelCallback
	asks          *Page
	bids          *Page
	queue         *cache.Queue
	marketService core.MarketService
}

func NewBook(ctx context.Context,
	market string,
	marketService core.MarketService,
	transact TransactCallback,
	cancel CancelCallback) *Book {
	return &Book{
		market:        market,
		createIndex:   make(map[string]bool),
		cancelIndex:   make(map[string]bool),
		transact:      transact,
		cancel:        cancel,
		asks:          NewPage(PageSideAsk),
		bids:          NewPage(PageSideBid),
		queue:         cache.NewQueue(ctx, market),
		marketService: marketService,
	}
}

func (book *Book) process(ctx context.Context, taker, maker *core.EngineOrder) (string, number.Integer, number.Integer) {
	taker.Assert()
	maker.Assert()

	matchedPrice := maker.Price
	makerAmount := maker.RemainingAmount
	makerFunds := makerAmount.Mul(matchedPrice)
	if maker.Side == PageSideBid {
		makerAmount = maker.RemainingFunds.Div(matchedPrice)
		makerFunds = maker.RemainingFunds
	}
	takerAmount := taker.RemainingAmount
	takerFunds := takerAmount.Mul(matchedPrice)
	if taker.Side == PageSideBid {
		takerAmount = taker.RemainingFunds.Div(matchedPrice)
		takerFunds = taker.RemainingFunds
	}
	matchedAmount, matchedFunds := makerAmount, makerFunds
	if takerAmount.Cmp(matchedAmount) < 0 || takerFunds.Cmp(matchedFunds) < 0 {
		matchedAmount, matchedFunds = takerAmount, takerFunds
	}

	taker.FilledAmount = taker.FilledAmount.Add(matchedAmount)
	taker.FilledFunds = taker.FilledFunds.Add(matchedFunds)
	if taker.Side == PageSideAsk {
		taker.RemainingAmount = taker.RemainingAmount.Sub(matchedAmount)
	}
	if taker.Side == PageSideBid {
		taker.RemainingFunds = taker.RemainingFunds.Sub(matchedFunds)
	}

	maker.FilledAmount = maker.FilledAmount.Add(matchedAmount)
	maker.FilledFunds = maker.FilledFunds.Add(matchedFunds)
	if maker.Side == PageSideAsk {
		maker.RemainingAmount = maker.RemainingAmount.Sub(matchedAmount)
	}
	if maker.Side == PageSideBid {
		maker.RemainingFunds = maker.RemainingFunds.Sub(matchedFunds)
	}

	tradeId := book.transact(taker, maker, matchedAmount)
	return tradeId, matchedAmount, matchedFunds
}

func (book *Book) processV2(ctx context.Context, taker, maker *core.EngineOrder) (string, number.Integer, number.Integer) {
	taker.Assert()
	maker.Assert()

	matchedPrice := maker.Price
	makerAmount := maker.RemainingAmount
	makerFunds := makerAmount.Mul(matchedPrice)
	//if maker.Side == PageSideBid {
	//	makerAmount = maker.RemainingFunds.Div(matchedPrice)
	//	makerFunds = maker.RemainingFunds
	//}
	takerAmount := taker.RemainingAmount
	takerFunds := takerAmount.Mul(matchedPrice)
	//if taker.Side == PageSideBid {
	//	takerAmount = taker.RemainingFunds.Div(matchedPrice)
	//	takerFunds = taker.RemainingFunds
	//}
	matchedAmount, matchedFunds := makerAmount, makerFunds
	if takerAmount.Cmp(matchedAmount) < 0 || takerFunds.Cmp(matchedFunds) < 0 {
		matchedAmount, matchedFunds = takerAmount, takerFunds
	}

	taker.FilledAmount = taker.FilledAmount.Add(matchedAmount)
	taker.FilledFunds = taker.FilledFunds.Add(matchedFunds)
	if taker.Side == PageSideAsk {
		taker.RemainingAmount = taker.RemainingAmount.Sub(matchedAmount)
	}
	if taker.Side == PageSideBid {
		taker.RemainingAmount = taker.RemainingAmount.Sub(matchedAmount)
		taker.RemainingFunds = taker.RemainingFunds.Sub(matchedFunds)
	}

	maker.FilledAmount = maker.FilledAmount.Add(matchedAmount)
	maker.FilledFunds = maker.FilledFunds.Add(matchedFunds)
	if maker.Side == PageSideAsk {
		maker.RemainingAmount = maker.RemainingAmount.Sub(matchedAmount)
	}
	if maker.Side == PageSideBid {
		maker.RemainingAmount = maker.RemainingAmount.Sub(matchedAmount)
		maker.RemainingFunds = maker.RemainingFunds.Sub(matchedFunds)
	}

	util.LogGoNumberIntegers(matchedAmount, matchedFunds)
	log.Println("taker")
	util.LogGoNumberIntegers(taker.RemainingAmount, taker.RemainingFunds, taker.FilledAmount, taker.FilledFunds)
	log.Println("maker")
	util.LogGoNumberIntegers(maker.RemainingAmount, maker.RemainingFunds, maker.FilledAmount, maker.FilledFunds)
	tradeId := book.transact(taker, maker, matchedAmount)
	return tradeId, matchedAmount, matchedFunds
}

func (book *Book) CreateOrder(ctx context.Context, order *core.EngineOrder) {
	if _, found := book.createIndex[order.Id]; found {
		return
	}
	book.createIndex[order.Id] = true

	if order.Side == PageSideAsk {
		opponents := make([]*core.EngineOrder, 0)
		book.bids.Iterate(func(opponent *core.EngineOrder) (number.Integer, number.Integer, bool) {
			if order.Type == core.OrderTypeLimit && opponent.Price.Cmp(order.Price) < 0 {
				return order.RemainingAmount.Zero(), order.RemainingFunds.Zero(), true
			}
			tradeId, matchedAmount, matchedFunds := book.processV2(ctx, order, opponent)
			book.cacheOrderEvent(ctx, cache.EventTypeOrderMatch, opponent.Side, opponent.Price, matchedAmount, matchedFunds, tradeId, opponent.Id, order.Id)
			opponents = append(opponents, opponent)
			return matchedAmount, matchedFunds, order.Filled()
		})
		for _, o := range opponents {
			if o.Filled() {
				book.bids.Remove(o)
			}
		}
		if !order.Filled() {
			if order.Type == core.OrderTypeLimit {
				book.asks.Put(order)
				book.cacheOrderEvent(ctx, cache.EventTypeOrderOpen, order.Side, order.Price, order.RemainingAmount, order.RemainingFunds, order.Id)
			} else {
				book.cancel(order, order.CreatedAt)
			}
		}
	} else if order.Side == PageSideBid {
		opponents := make([]*core.EngineOrder, 0)
		book.asks.Iterate(func(opponent *core.EngineOrder) (number.Integer, number.Integer, bool) {
			if order.Type == core.OrderTypeLimit && opponent.Price.Cmp(order.Price) > 0 {
				return order.RemainingAmount.Zero(), order.RemainingFunds.Zero(), true
			}
			tradeId, matchedAmount, matchedFunds := book.processV2(ctx, order, opponent)
			book.cacheOrderEvent(ctx, cache.EventTypeOrderMatch, opponent.Side, opponent.Price, matchedAmount, matchedFunds, tradeId, opponent.Id, order.Id)
			opponents = append(opponents, opponent)
			return matchedAmount, matchedFunds, order.Filled()
		})
		for _, o := range opponents {
			if o.Filled() {
				book.asks.Remove(o)
			}
		}
		if !order.Filled() {
			if order.Type == core.OrderTypeLimit {
				book.bids.Put(order)
				book.cacheOrderEvent(ctx, cache.EventTypeOrderOpen, order.Side, order.Price, order.RemainingAmount, order.RemainingFunds, order.Id)
			} else {
				book.cancel(order, order.CreatedAt)
			}
		}
	}
}

func (book *Book) CancelOrder(ctx context.Context, order *core.EngineOrder, cancelledAt time.Time) {
	if _, found := book.cancelIndex[order.Id]; found {
		return
	}
	book.cancelIndex[order.Id] = true

	var internalOrder *core.EngineOrder
	if order.Side == PageSideAsk {
		internalOrder = book.asks.Remove(order)
	} else if order.Side == PageSideBid {
		internalOrder = book.bids.Remove(order)
	} else {
		log.Panicln(order)
	}
	if internalOrder != nil {
		book.cancel(internalOrder, cancelledAt)
		book.cacheOrderEvent(ctx, cache.EventTypeOrderCancel, order.Side, order.Price, order.RemainingAmount, order.RemainingFunds, order.Id)
	} else {
		book.cancel(order, cancelledAt)
	}
}

func (book *Book) Run(ctx context.Context) {
	go book.queue.Loop(ctx)

	fullCacheTicker := time.NewTicker(time.Second * 30)
	defer fullCacheTicker.Stop()

	bestCacheTicker := time.NewTicker(time.Second * 2)
	defer bestCacheTicker.Stop()

	book.cacheList(ctx, 0)

	for {
		select {
		//case event := <-book.events:
		//	if event.Action.Action == OrderActionCreate {
		//		book.createOrder(ctx, event.Order)
		//	} else if event.Action.Action == OrderActionCancel {
		//		book.cancelOrder(ctx, event.Order, event.Action.CreatedAt)
		//	} else {
		//		log.Panicln(event)
		//	}
		case <-fullCacheTicker.C:
			book.cacheList(ctx, 0)
		case <-bestCacheTicker.C:
			book.cacheList(ctx, 1)
		}
	}
}

func (book *Book) cacheList(ctx context.Context, limit int) {
	//event := fmt.Sprintf("BOOK-T%d", limit)
	asks, bids := book.asks.List(limit, true), book.bids.List(limit, true)
	//data := map[string]interface{}{
	//	"asks": asks,
	//	"bids": bids,
	//}
	go func() {
		bidss, _ := json.Marshal(bids)
		askss, _ := json.Marshal(asks)
		_ = book.marketService.UpdateOptionMarket(book.market, string(bidss), string(askss))
	}()
	//book.queue.AttachEvent(ctx, event, data)
}

func (book *Book) cacheOrderEvent(ctx context.Context, event, side string, price, amount, funds number.Integer, tradeAndOrderIds ...string) {
	if amount.IsZero() {
		amount = funds.Div(price)
	} else if funds.IsZero() {
		funds = price.Mul(amount)
	}
	data := map[string]interface{}{
		"side":   side,
		"price":  price,
		"amount": amount,
		"funds":  funds,
	}

	switch event {
	case cache.EventTypeOrderOpen, cache.EventTypeOrderCancel: // order open or cancel event
		data["order_id"] = tradeAndOrderIds[0]
	case cache.EventTypeOrderMatch: // order match event
		data["trade_id"] = tradeAndOrderIds[0]
		data["maker_id"] = tradeAndOrderIds[1]
		data["taker_id"] = tradeAndOrderIds[2]
	}

	//book.queue.AttachEvent(ctx, event, data)
}

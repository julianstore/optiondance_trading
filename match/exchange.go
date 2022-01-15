package match

import (
	"context"
	"github.com/MixinNetwork/go-number"
	"github.com/gofrs/uuid"
	"github.com/ugorji/go/codec"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"math/big"
	"option-dance/core"
	"option-dance/match/engine"
	"option-dance/pkg/log"
	"sync"
	"time"
)

const (
	PollInterval   = 100 * time.Millisecond
	EventQueueSize = 8192
)

type Job interface {
	Run(ctx context.Context)
}

type Exchange struct {
	books          map[string]*engine.Book
	booksMutex     *sync.RWMutex
	codec          codec.Handle
	snapshots      map[string]bool
	mutexes        *tmap
	events         chan *engine.OrderEvent
	propertyStore  core.PropertyStore
	marketService  core.MarketService
	messageStore   core.MessageStore
	messageService core.MessageService
	messageBuilder core.MessageBuilder
	tradeService   core.TradeService
	orderService   core.OrderService
	transferStore  core.TransferStore
	rawTxStore     core.RawTxStore
	utxoStore      core.UtxoStore
	utxoDispatcher UtxoDispatcher
	utxoSyncer     UtxoSyncer
	spentSyncer    SpentSyncer
	multiSigner    MultiSigner
	notifier       core.Notifier
}

func NewExchange(
	propertyStore core.PropertyStore,
	marketService core.MarketService,
	messageStore core.MessageStore,
	messageService core.MessageService,
	messageBuilder core.MessageBuilder,
	orderService core.OrderService,
	transferStore core.TransferStore,
	utxoStore core.UtxoStore,
	rawTxStore core.RawTxStore,
	utxoDispatcher UtxoDispatcher,
	utxoSycner UtxoSyncer,
	spentSyncer SpentSyncer,
	multiSigner MultiSigner,
	tradeService core.TradeService,
	notifier core.Notifier,

) *Exchange {
	return &Exchange{
		codec:          new(codec.MsgpackHandle),
		books:          make(map[string]*engine.Book),
		booksMutex:     &sync.RWMutex{},
		snapshots:      make(map[string]bool),
		mutexes:        newTmap(),
		events:         make(chan *engine.OrderEvent, EventQueueSize),
		propertyStore:  propertyStore,
		marketService:  marketService,
		messageStore:   messageStore,
		orderService:   orderService,
		messageService: messageService,
		messageBuilder: messageBuilder,
		transferStore:  transferStore,
		utxoStore:      utxoStore,
		rawTxStore:     rawTxStore,
		utxoDispatcher: utxoDispatcher,
		utxoSyncer:     utxoSycner,
		spentSyncer:    spentSyncer,
		tradeService:   tradeService,
		multiSigner:    multiSigner,
		notifier:       notifier,
	}
}

func (ex *Exchange) Run(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("ExchangeRun-Panic", zap.Any("err", err))
		}
	}()

	g := new(errgroup.Group)

	g.Go(func() error { return ex.handleOrderEvent(ctx) })       //handle order event to order book
	g.Go(func() error { return ex.PollOrderActions(ctx) })       //parse orderActions to match engine
	g.Go(func() error { return ex.LoopingSendTransaction(ctx) }) //send raw transaction to mixin network
	g.Go(func() error { return ex.utxoSyncer.Run(ctx) })         //sync utxos from mixin network
	g.Go(func() error { return ex.utxoDispatcher.Run(ctx) })     //parse utxo memo to biz actions
	g.Go(func() error { return ex.multiSigner.Run(ctx) })        //handler transfer and multisig
	g.Go(func() error { return ex.spentSyncer.Run(ctx) })        //spent syncer
	g.Go(func() error { return ex.LoopingSendMessage(ctx) })     //send message notify

	zap.L().Info("engine started.")
	if err := g.Wait(); err != nil {
		zap.L().Error("Looping Error", zap.Error(err))
		return err
	}
	return nil
}

func (ex *Exchange) buildBook(ctx context.Context, market string) *engine.Book {
	return engine.NewBook(ctx, market, ex.marketService,
		func(taker, maker *core.EngineOrder, amount number.Integer) string {
			for {
				tradeId, err := ex.tradeService.Transact(ctx, taker, maker, amount)
				if err == nil {
					return tradeId
				}
				log.ErrorInfo("Engine Transact CALLBACK", err)
				time.Sleep(PollInterval)
			}
		},
		func(order *core.EngineOrder, cancelledAt time.Time) {
			for {
				err := ex.tradeService.CancelOrder(ctx, order, cancelledAt)
				if err == nil {
					break
				}
				log.ErrorInfo("Engine Cancel CALLBACK", err)
				time.Sleep(PollInterval)
			}
		},
	)
}

func (ex *Exchange) AttachOrderEvent(ctx context.Context, order *core.EngineOrder, action *core.Action) {
	if order.Side != engine.PageSideAsk && order.Side != engine.PageSideBid {
		zap.L().Panic("order", zap.Any("order", order), zap.Any("action", action))
	}
	if order.Type != core.OrderTypeLimit && order.Type != core.OrderTypeMarket {
		zap.L().Panic("order", zap.Any("order", order), zap.Any("action", action))
	}
	if action.Action != engine.OrderActionCancel {
		order.Assert()
	}
	ex.events <- &engine.OrderEvent{Order: order, Action: action}
}

func (ex *Exchange) handleOrderEvent(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-ex.events:
			ex.booksMutex.Lock()
			if event.Action.Action == engine.OrderActionCreate {
				ex.books[event.Order.InstrumentName].CreateOrder(ctx, event.Order)
				go func() {
					err := ex.notifier.CreateOrderGroupNotify(ctx, *event.Order)
					if err != nil {
						zap.L().Error("sendCreateOrderGroupNotify", zap.Error(err))
					}
				}()
			} else if event.Action.Action == engine.OrderActionCancel {
				ex.books[event.Order.InstrumentName].CancelOrder(ctx, event.Order, event.Action.CreatedAt)
			} else {
				zap.L().Error("event error", zap.Any("event", event))
			}
			ex.booksMutex.Unlock()
		}
	}
}

type tmap struct {
	sync.Map
}

func newTmap() *tmap {
	return &tmap{
		Map: sync.Map{},
	}
}

func (m *tmap) fetch(user, asset string) *sync.Mutex {
	uu, err := uuid.FromString(user)
	if err != nil {
		panic(user)
	}
	u := new(big.Int).SetBytes(uu.Bytes())
	au, err := uuid.FromString(asset)
	if err != nil {
		panic(asset)
	}
	a := new(big.Int).SetBytes(au.Bytes())
	s := new(big.Int).Add(u, a)
	key := new(big.Int).Mod(s, big.NewInt(100)).String()
	if _, found := m.Load(key); !found {
		m.Store(key, new(sync.Mutex))
	}
	val, _ := m.Load(key)
	return val.(*sync.Mutex)
}

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Snapshot struct {
	SnapshotId string `json:"snapshot_id"`
	Amount     string `json:"amount"`
	Asset      struct {
		AssetId string `json:"asset_id"`
	} `json:"asset"`
	CreatedAt time.Time `json:"created_at"`

	TraceId    string `json:"trace_id"`
	UserId     string `json:"user_id"`
	OpponentId string `json:"opponent_id"`
	Data       string `json:"data"`
}

func QuotePrecision(assetId string) uint8 {
	return 8
}

func QuoteMinimum(assetId string) number.Decimal {
	return number.Zero()
}

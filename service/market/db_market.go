package market

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/rpc/deribit"
	"time"
)

var (
	AskPriceRatio = decimal.NewFromFloat(1.1)
	BidPriceRatio = decimal.NewFromFloat(0.9)
	MaxDiffPrice  = decimal.NewFromFloat(150)
)

type DbMarketService struct {
	dbMarketTrackerStore core.DbMarketTrackerStore
	dbMarketStore        core.DbMarketStore
	dbClient             *deribit.Client
	propertyStore        core.PropertyStore
}

type IntervalConfig struct {
	SyncMarketInterval     int64
	SyncIndexPriceInterval int64
}

func NewDbMarketService(
	dbMarketTrackerStore core.DbMarketTrackerStore,
	dbMarketStore core.DbMarketStore,
	propertyStore core.PropertyStore,
) core.DbMarketService {
	dbClient := deribit.C()
	return &DbMarketService{
		dbMarketTrackerStore: dbMarketTrackerStore,
		dbMarketStore:        dbMarketStore,
		dbClient:             dbClient,
		propertyStore:        propertyStore,
	}
}

func (s *DbMarketService) getIntervalConfig() (c IntervalConfig) {
	d := config.Cfg.Deribit
	c = IntervalConfig{
		SyncMarketInterval:     d.SyncMarketInterval,
		SyncIndexPriceInterval: d.SyncIndexPriceInterval,
	}
	if d.SyncIndexPriceInterval <= 0 {
		c.SyncIndexPriceInterval = 5
	}
	if d.SyncMarketInterval <= 0 {
		c.SyncMarketInterval = 2
	}
	return
}

func (s *DbMarketService) SyncMarket(ctx context.Context) error {

	var syncMarketInterval = time.Millisecond
	c := s.getIntervalConfig()
	for {
		select {
		case <-ctx.Done():
			zap.L().Error("sync market error", zap.Error(ctx.Err()))
			return ctx.Err()
		case <-time.After(syncMarketInterval):
			if err := s.RunSync(ctx); err != nil {
				syncMarketInterval = time.Duration(c.SyncMarketInterval) * time.Second
				zap.L().Error("sync market error", zap.Error(err))
			} else {
				syncMarketInterval = time.Duration(c.SyncMarketInterval) * time.Second
			}
		}
	}
}

func (s *DbMarketService) SyncIndexPrice(ctx context.Context) error {

	var syncIndexPriceInterval = time.Millisecond
	c := s.getIntervalConfig()
	for {
		select {
		case <-ctx.Done():
			zap.L().Error("SyncIndexPrice error", zap.Error(ctx.Err()))
			return ctx.Err()
		case <-time.After(syncIndexPriceInterval):
			_ = s.syncIndexPrice(ctx)
			syncIndexPriceInterval = time.Duration(c.SyncIndexPriceInterval) * time.Second
		}
	}
}

func (s *DbMarketService) RunSync(ctx context.Context) error {

	trackerList, err := s.dbMarketTrackerStore.ListAll(ctx)
	if err != nil {
		return err
	}
	trackers, err := s.dbMarketTrackerStore.ToMarketTracker(ctx, trackerList)
	if err != nil {
		return err
	}

	for _, e := range trackers {

		if e.Status == core.TrackStatusOn {
			res, err := s.dbClient.GetOrderBook(e.DbInstrumentName, 1)
			if err != nil {
				return err
			}

			book := res.Result
			e.Bids, e.Asks = "[]", "[]"
			if e.AskTrack {
				if len(book.Asks) > 0 {
					dbOrderPrice := book.Asks[0][0]
					dbOrderUSDPrice := s.CalcOrderPrice(ctx, core.PageSideAsk, book.IndexPrice, dbOrderPrice)
					dbOrderSize := decimal.NewFromFloat(book.Asks[0][1])
					e.Asks, err = s.GetOdOrderBookEntry(ctx, dbOrderUSDPrice, dbOrderSize, core.PageSideAsk)
					if err != nil {
						return err
					}
				}
			}
			if e.BidTrack {
				if len(book.Bids) > 0 {
					dbOrderPrice := book.Bids[0][0]
					dbOrderUSDPrice := s.CalcOrderPrice(ctx, core.PageSideBid, book.IndexPrice, dbOrderPrice)
					dbOrderSize := decimal.NewFromFloat(book.Bids[0][1])
					e.Bids, err = s.GetOdOrderBookEntry(ctx, dbOrderUSDPrice, dbOrderSize, core.PageSideBid)
					if err != nil {
						return err
					}
				}
			}
			if err = s.SaveTrackerToOdMarket(ctx, e); err != nil {
				return err
			}
		} else {
			if err = s.dbMarketStore.CloseMarket(ctx, e.OdInstrumentName); err != nil {
				return err
			}
		}
	}

	// close market that no more exists
	all, err := s.dbMarketStore.ListAll(ctx)
	if err != nil {
		return err
	}

	var trackerNames []string
	for _, e := range trackers {
		trackerNames = append(trackerNames, e.OdInstrumentName)
	}

	for _, e := range all {
		if !govalidator.IsIn(e.InstrumentName, trackerNames...) {
			if err = s.dbMarketStore.CloseMarket(ctx, e.InstrumentName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *DbMarketService) syncIndexPrice(ctx context.Context) error {
	price, err := s.dbClient.GetIndexPrice("btc_usd")
	if err != nil {
		return err
	}
	indexPrice := price.Result.IndexPrice
	v := decimal.NewFromFloat(indexPrice).String()
	return s.propertyStore.WriteProperty(ctx, core.DbBtcUsdIndexPrice, v)
}

func (s *DbMarketService) GetOdOrderBookEntry(ctx context.Context, usdPrice, amount decimal.Decimal, side string) (string, error) {
	entry := core.Entry{
		Side:   side,
		Price:  usdPrice.StringFixed(2),
		Amount: amount.StringFixedBank(2),
		Funds:  "",
	}
	marshal, err := json.Marshal([]core.Entry{entry})
	if err != nil {
		return "[]", err
	}
	return string(marshal), nil
}

func (s *DbMarketService) CalcOrderPrice(ctx context.Context, side string, indexPrice, dbOrderPrice float64) (orderUsdPrice decimal.Decimal) {
	orderBtcPrice := decimal.NewFromFloat(dbOrderPrice)
	dbUsdPrice := decimal.NewFromFloat(indexPrice).Mul(orderBtcPrice)
	x := dbUsdPrice.Mul(decimal.NewFromFloat(0.1))
	var reachMaxDiff = x.Cmp(MaxDiffPrice) >= 0
	if side == core.PageSideAsk {
		orderUsdPrice = dbUsdPrice.Mul(AskPriceRatio)
		if reachMaxDiff {
			orderUsdPrice = dbUsdPrice.Add(MaxDiffPrice)
		}
	}
	if side == core.PageSideBid {
		orderUsdPrice = dbUsdPrice.Mul(BidPriceRatio)
		if reachMaxDiff {
			orderUsdPrice = dbUsdPrice.Sub(MaxDiffPrice)
		}
	}
	return orderUsdPrice
}

func (s *DbMarketService) SaveTrackerToOdMarket(ctx context.Context, tracker *core.MarketTracker) error {
	i, err := core.ParseInstrument(tracker.OdInstrumentName)
	if err != nil {
		return err
	}

	rowAffected, err := s.dbMarketStore.UpdateBidAsks(tracker.OdInstrumentName, tracker.Bids, tracker.Asks)
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		dbMarket := core.DbMarket{
			InstrumentName:      tracker.OdInstrumentName,
			StrikePrice:         i.StrikePrice,
			ExpirationDate:      i.ExpirationDate,
			ExpirationDateStr:   i.ExpirationDate.Format("2006-01-02"),
			ExpirationTimestamp: i.ExpirationTimestamp,
			OptionType:          i.OptionType,
			DeliveryType:        i.DeliveryType,
			QuoteCurrency:       i.QuoteCurrency,
			BaseCurrency:        i.BaseCurrency,
			Bids:                tracker.Bids,
			Asks:                tracker.Asks,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			Status:              0,
		}
		if err := s.dbMarketStore.Save(ctx, dbMarket); err != nil {
			return err
		}
	}
	return nil
}

func (s *DbMarketService) ListExpiryDatesByPrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []core.OptionMarketDTO, err error) {
	var list []core.OptionMarketDTO
	list, err = s.dbMarketStore.ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil, err
	}
	for _, dto := range list {
		optionMarketList = append(optionMarketList, dto)
	}
	return
}

func (s *DbMarketService) ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int64) {
	optionMarkets, err := s.dbMarketStore.ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil
	}
	for _, e := range optionMarkets {
		prices = append(prices, e.StrikePrice)
	}
	return
}

func (s *DbMarketService) ListMarketsByInstrument(name string) (d core.OptionMarketDTO, err error) {
	d, err = s.dbMarketStore.FindByInstrumentName(name)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (s *DbMarketService) ListMarketsByInstrumentOdAndDb(name string) (d core.OptionMarketDTO, err error) {
	return d, nil
}

func (s *DbMarketService) UpdateOptionMarket(instrumentName string, bids, asks string) (err error) {
	return nil
}

func (s *DbMarketService) ListExpiryDatesByPriceOdAndDb(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []core.OptionMarketDTO, err error) {
	return nil, nil
}

func (s *DbMarketService) ListMarketStrikePriceOdAndDb(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int) {
	return nil
}

func (s *DbMarketService) ListDates(ctx context.Context, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error) {
	return nil, err
}

func (s *DbMarketService) ListByDate(ctx context.Context, date string) (dateList []core.FullMarketByDateDTO, err error) {
	return nil, err
}

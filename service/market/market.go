package market

import (
	"context"
	"option-dance/core"
	"sort"
	"time"
)

func New(marketStore core.MarketStore, dbMarketStore core.DbMarketStore) core.MarketService {
	return &service{marketStore: marketStore, dbMarketStore: dbMarketStore}
}

type service struct {
	marketStore   core.MarketStore
	dbMarketStore core.DbMarketStore
}

func (s *service) ListExpiryDatesByPrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []core.OptionMarketDTO, err error) {
	var list []core.OptionMarketDTO
	list, err = s.marketStore.ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil, err
	}
	for _, dto := range list {
		optionMarketList = append(optionMarketList, dto)
	}
	return
}

func (s *service) ListExpiryDatesByPriceOdAndDb(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarketList []core.OptionMarketDTO, err error) {
	var marketMap = make(map[string]bool, 0)
	odMarkets, err := s.marketStore.ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil, err
	}
	for _, dto := range odMarkets {
		marketMap[dto.InstrumentName] = true
		optionMarketList = append(optionMarketList, dto)
	}

	dbMarkets, err := s.dbMarketStore.ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil, err
	}
	for _, dto := range dbMarkets {
		if _, found := marketMap[dto.InstrumentName]; !found {
			optionMarketList = append(optionMarketList, dto)
			marketMap[dto.InstrumentName] = true
		}
	}
	return
}

func (s *service) ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int64) {
	optionMarkets, err := s.marketStore.ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil
	}
	for _, e := range optionMarkets {
		prices = append(prices, e.StrikePrice)
	}
	return
}

func (s *service) ListMarketStrikePriceOdAndDb(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (prices []int) {
	odMarkets, err := s.marketStore.ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return nil
	}
	var (
		priceMap = make(map[int64]bool, 0)
	)
	for _, e := range odMarkets {
		priceMap[e.StrikePrice] = true
		prices = append(prices, int(e.StrikePrice))
	}

	dbMarkets, err := s.dbMarketStore.ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency)
	if err != nil {
		return prices
	}
	for _, e := range dbMarkets {
		if _, found := priceMap[e.StrikePrice]; !found {
			prices = append(prices, int(e.StrikePrice))
			priceMap[e.StrikePrice] = true
		}
	}
	sort.Ints(prices)
	return
}

func (s *service) ListMarketsByInstrument(name string) (d core.OptionMarketDTO, err error) {
	d, err = s.marketStore.FindByInstrumentName(name)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (s *service) ListMarketsByInstrumentOdAndDb(name string) (d core.OptionMarketDTO, err error) {
	d, err = s.marketStore.FindByInstrumentName(name)
	if err != nil {
		return d, err
	}
	if d.InstrumentName == "" || core.Expired(time.Now().Unix(), d.ExpirationTimestamp) || (len(d.AskList) == 0 && len(d.BidList) == 0) {
		d, err = s.dbMarketStore.FindByInstrumentName(name)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func (s *service) UpdateOptionMarket(instrumentName string, bids, asks string) error {

	rowsAffected, err := s.marketStore.UpdateBidAsks(instrumentName, bids, asks)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		i, err := core.ParseInstrument(instrumentName)
		if err != nil {
			return err
		}
		market := core.OptionMarket{
			InstrumentName:      instrumentName,
			StrikePrice:         i.StrikePrice,
			ExpirationDate:      i.ExpirationDate,
			ExpirationDateStr:   i.ExpirationDate.Format("2006-01-02"),
			ExpirationTimestamp: i.ExpirationTimestamp,
			OptionType:          i.OptionType,
			DeliveryType:        i.DeliveryType,
			QuoteCurrency:       i.QuoteCurrency,
			BaseCurrency:        i.BaseCurrency,
			Bids:                bids,
			Asks:                asks,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			Status:              0,
		}
		err = s.marketStore.Create(market)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ListDates(ctx context.Context, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error) {
	return s.marketStore.ListDates(ctx, deliveryType, quoteCurrency, baseCurrency)
}

func (s *service) ListByDate(ctx context.Context, date string) (dateList []core.FullMarketByDateDTO, err error) {
	return s.marketStore.ListByDate(ctx, date)
}

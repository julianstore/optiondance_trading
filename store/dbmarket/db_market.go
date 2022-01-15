package dbmarket

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"option-dance/core"
	store2 "option-dance/store"
	"time"
)

type dbMarketStore struct {
	db *gorm.DB
}

func NewDbMarketStore(db *gorm.DB) core.DbMarketStore {
	return &dbMarketStore{db: db}
}

func (s *dbMarketStore) Save(ctx context.Context, market core.DbMarket) error {
	return s.db.Model(core.DbMarket{}).Create(&market).Error
}

func (s *dbMarketStore) UpdateBidAsks(instrumentName string, bids, asks string) (rowsAffected int64, err error) {
	tx := s.db.Model(core.DbMarket{}).Where("instrument_name = ?", instrumentName).Updates(map[string]interface{}{
		"bids":       bids,
		"asks":       asks,
		"updated_at": time.Now(),
		"status":     core.TrackStatusOn,
	})
	return tx.RowsAffected, nil
}

func (s *dbMarketStore) FindByInstrumentName(instrumentName string) (o core.OptionMarketDTO, err error) {
	var market core.OptionMarket
	if err = s.db.Model(core.DbMarket{}).
		Where("instrument_name=? and status = 0", instrumentName).Find(&market).Error; store2.CheckErr(err) != nil {
		return core.OptionMarketDTO{}, err
	}
	if market.InstrumentName == "" {
		return o, nil
	}
	o, err = core.ConvertToMarketDTO(market)
	if err != nil {
		return core.OptionMarketDTO{}, err
	}
	return
}

func (s *dbMarketStore) ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (list []core.OptionMarketDTO, err error) {
	var markets []core.OptionMarket
	now := time.Now().Unix()

	var sideField = "asks"
	if side == core.PageSideBid {
		sideField = "bids"
	}

	condition := fmt.Sprintf(`
%s != '[]' 
and strike_price = '%s' 
and option_type = '%s' 
and delivery_type = '%s' 
and quote_currency = '%s' 
and base_currency = '%s' 
and expiration_timestamp > %d and status = 0`, sideField, price, optionType, deliveryType, quoteCurrency, baseCurrency, now)

	if err := s.db.Model(core.DbMarket{}).Where(condition).
		Order("expiration_timestamp asc").
		Find(&markets).Error; store2.CheckErr(err) != nil {
		return list, err
	}
	for _, e := range markets {
		dto, err := core.ConvertToMarketDTO(e)
		if err != nil {
			return nil, err
		}
		list = append(list, dto)
	}
	return list, nil
}

func (s *dbMarketStore) ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarkets []core.OptionMarketDTO, err error) {
	var sideCondition string
	now := time.Now().Unix()
	if side == "BID" {
		sideCondition = fmt.Sprintf("asks != '[]' and option_type = '%s' and delivery_type = '%s' "+
			"and quote_currency = '%s' and base_currency = '%s' and expiration_timestamp > %d and status = 0",
			optionType, deliveryType, quoteCurrency, baseCurrency, now)
	} else {
		sideCondition = fmt.Sprintf("bids != '[]' and option_type = '%s' and delivery_type = '%s' "+
			"and quote_currency = '%s' and base_currency = '%s' and expiration_timestamp > %d and status = 0",
			optionType, deliveryType, quoteCurrency, baseCurrency, now)
	}
	s.db.Model(core.OptionMarket{}).
		Raw(fmt.Sprintf("SELECT strike_price FROM db_market where %s group by strike_price order by strike_price", sideCondition)).
		Scan(&optionMarkets)
	return
}

func (s *dbMarketStore) DeleteByInstrumentName(ctx context.Context, instrumentName string) error {
	return nil
}

func (s *dbMarketStore) Create(market core.OptionMarket) error {
	return nil
}

func (s *dbMarketStore) CloseMarket(ctx context.Context, name string) error {
	return s.db.Model(core.DbMarket{}).Where("instrument_name = ?", name).Updates(map[string]interface{}{
		"status":     core.TrackStatusOff,
		"bids":       "[]",
		"asks":       "[]",
		"updated_at": time.Now(),
	}).Error
}

func (s *dbMarketStore) ListAll(ctx context.Context) (markets []*core.DbMarket, err error) {
	if err = s.db.Model(core.DbMarket{}).WithContext(ctx).Where("status = ?", 0).Find(&markets).Error; err != nil {
		return nil, err
	}
	return
}

func (s *dbMarketStore) ListDates(ctx context.Context, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error) {
	return nil, err
}

func (s *dbMarketStore) ListByDate(ctx context.Context, date string) (dateList []core.FullMarketByDateDTO, err error) {
	return nil, err
}

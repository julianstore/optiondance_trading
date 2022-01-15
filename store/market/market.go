package market

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"gorm.io/gorm"
	"option-dance/core"
	store2 "option-dance/store"
	"time"
)

func New(db *gorm.DB) core.MarketStore {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (s *store) Create(market core.OptionMarket) error {
	err := s.db.Model(core.OptionMarket{}).Create(&market).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *store) FindByInstrumentName(instrumentName string) (o core.OptionMarketDTO, err error) {
	var market core.OptionMarket
	if err = s.db.Model(core.OptionMarket{}).
		Where("instrument_name=?", instrumentName).Find(&market).Error; store2.CheckErr(err) != nil {
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

func (s *store) ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (list []core.OptionMarketDTO, err error) {

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
and expiration_timestamp > %d`, sideField, price, optionType, deliveryType, quoteCurrency, baseCurrency, now)

	if err := s.db.Model(core.OptionMarket{}).Where(condition).
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

func (s *store) ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarkets []core.OptionMarketDTO, err error) {
	var sideCondition string
	now := time.Now().Unix()
	if side == "BID" {
		sideCondition = fmt.Sprintf("asks != '[]' and option_type = '%s' and delivery_type = '%s' "+
			"and quote_currency = '%s' and base_currency = '%s' and expiration_timestamp > %d",
			optionType, deliveryType, quoteCurrency, baseCurrency, now)
	} else {
		sideCondition = fmt.Sprintf("bids != '[]' and option_type = '%s' and delivery_type = '%s' "+
			"and quote_currency = '%s' and base_currency = '%s' and expiration_timestamp > %d",
			optionType, deliveryType, quoteCurrency, baseCurrency, now)
	}
	s.db.Model(core.OptionMarket{}).
		Raw(fmt.Sprintf("SELECT strike_price FROM option_market where %s group by strike_price order by strike_price", sideCondition)).
		Scan(&optionMarkets)
	return
}

func (s *store) UpdateBidAsks(instrumentName string, bids, asks string) (rowsAffected int64, err error) {
	tx := s.db.Model(core.OptionMarket{}).Where("instrument_name = ?", instrumentName).
		Updates(map[string]interface{}{
			"bids":       bids,
			"asks":       asks,
			"updated_at": time.Now(),
		})
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func (s *store) DeleteByInstrumentName(ctx context.Context, instrumentName string) error {
	if err := s.db.Model(core.OptionMarket{}).Where("instrument_name=?", instrumentName).Delete(core.OptionMarket{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) ListDates(ctx context.Context, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error) {

	now := time.Now().Unix()
	condition := fmt.Sprintf("delivery_type = '%s' "+
		"and quote_currency = '%s' and base_currency = '%s' and expiration_timestamp > %d", deliveryType, quoteCurrency, baseCurrency, now)

	if err := s.db.Model(core.OptionMarket{}).
		Raw(fmt.Sprintf("SELECT distinct expiration_date_str FROM option_market where %s group by strike_price order by strike_price", condition)).
		Scan(&dateList).Error; err != nil {
		return nil, err
	}
	return
}

func (s *store) ListByDate(ctx context.Context, date string) (list []core.FullMarketByDateDTO, err error) {
	var marketList []core.OptionMarket
	if err := s.db.Model(core.OptionMarket{}).Where("expiration_date_str = ?", date).
		Scan(&marketList).Error; err != nil {
		return nil, err
	}
	list, err = s.convertToFullMarketByDateDTO(marketList)
	if err != nil {
		return nil, err
	}
	return
}

func (s *store) convertToFullMarketByDateDTO(marketList []core.OptionMarket) (list []core.FullMarketByDateDTO, err error) {
	var resMap = treemap.NewWithIntComparator()
	for _, e := range marketList {
		dto, err := core.ConvertToMarketDTO(e)
		if err != nil {
			return nil, err
		}
		key := int(e.StrikePrice)

		fullMarket := core.FullMarketByDateDTO{StrikePrice: key}
		if dto.OptionType == core.OptionTypeCALL {
			fullMarket.Call = dto
		}
		if dto.OptionType == core.OptionTypePUT {
			fullMarket.Put = dto
		}
		if v, found := resMap.Get(key); found {
			fullMarket := v.(core.FullMarketByDateDTO)
			if dto.OptionType == core.OptionTypeCALL {
				fullMarket.Call = dto
			} else {
				fullMarket.Put = dto
			}
			resMap.Put(key, fullMarket)
		} else {
			resMap.Put(key, fullMarket)
		}
	}
	resMap.Each(func(key interface{}, value interface{}) {
		dateDTO := value.(core.FullMarketByDateDTO)
		list = append(list, dateDTO)
	})
	return
}

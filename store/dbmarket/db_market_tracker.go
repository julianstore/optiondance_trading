package dbmarket

import (
	"context"
	"gorm.io/gorm"
	"option-dance/core"
	"strings"
	"time"
)

type dbMarketTrackerStore struct {
	db *gorm.DB
}

func NewDbMarketTrackerStore(db *gorm.DB) core.DbMarketTrackerStore {
	return &dbMarketTrackerStore{db: db}
}

func (s *dbMarketTrackerStore) Save(ctx context.Context, date, marketType, strikePrices string) {
	panic("implement me")
}

func (s *dbMarketTrackerStore) ListAll(ctx context.Context) (trackers []*core.DbMarketTracker, err error) {
	if err = s.db.Model(core.DbMarketTracker{}).Find(&trackers).Error; err != nil {
		return nil, err
	}
	return
}

// ToMarketTracker  from db market tracker to MarketTrackers
func (s *dbMarketTrackerStore) ToMarketTracker(ctx context.Context, trackers []*core.DbMarketTracker) (marketTrackers []*core.MarketTracker, err error) {

	var trackerMap = make(map[string]*core.MarketTracker, 0)
	for _, e := range trackers {
		optionType := core.OptionTypePUT
		if strings.Contains(strings.ToLower(e.Type), "call") {
			optionType = core.OptionTypeCALL
		}
		expiry, err := time.Parse("2006/01/02", e.Date)
		if err != nil {
			return nil, err
		}
		split := strings.Split(e.StrikePrices, "/")
		for _, strikePrice := range split {

			var (
				dbInstrument       = core.ToDbInstrument(e.BaseCurrency, strikePrice, optionType, expiry)
				odInstrument       = core.ToOdInstrument(e.DeliveryType, "", e.BaseCurrency, strikePrice, optionType, expiry)
				key                = dbInstrument
				askTrack, bidTrack bool
				isAskSide          = strings.Contains(strings.ToLower(e.Type), "sell")
			)

			if e.Status == core.TrackStatusOn {
				askTrack = isAskSide
				bidTrack = !isAskSide
			}
			if v, ok := trackerMap[key]; !ok {
				trackerMap[key] = &core.MarketTracker{
					DbInstrumentName: dbInstrument,
					OdInstrumentName: odInstrument,
					AskTrack:         askTrack,
					BidTrack:         bidTrack,
					Status:           e.Status,
				}
			} else {
				if v.Status == core.TrackStatusOff {
					v.Status = e.Status
				}
				if e.Status == core.TrackStatusOn {
					if !v.AskTrack {
						v.AskTrack = isAskSide
					}
					if !v.BidTrack {
						v.BidTrack = !isAskSide
					}
				}
			}
		}
	}
	for _, v := range trackerMap {
		marketTrackers = append(marketTrackers, v)
	}
	return marketTrackers, nil
}

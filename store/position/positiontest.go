package position

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"option-dance/core"
	time2 "option-dance/pkg/time"
	"option-dance/pkg/util"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func NewPositionTestStore() core.PositionStore {
	return &positionTestStore{}
}

type positionTestStore struct {
}

func (p *positionTestStore) FindByPositionId(ctx context.Context, positionId string) (position *core.Position, err error) {
	return nil, nil
}

func (p *positionTestStore) FindByPositionIdAndUid(ctx context.Context, uid, positionId string) (position *core.Position, err error) {
	return nil, nil
}

func (p *positionTestStore) AutoExercisePositionWithName(ctx context.Context, instrumentName string, updateTime time.Time) (err error) {
	return nil
}

func (p *positionTestStore) ExercisePositionWithID(ctx context.Context, positionId, instrumentName string, updateTime time.Time) (err error) {
	return nil
}

func (p *positionTestStore) ListNamesByDateAndSide(ctx context.Context, time time.Time, status int64, side, deliveryType string) (names []string, err error) {
	return nil, nil
}

func (p *positionTestStore) loadTestPositions(dateString string) (positions []*core.Position, err error) {

	join := filepath.Join(util.AbsPath(), fmt.Sprintf("testdata/%s/options.json", dateString))
	bytes, err := ioutil.ReadFile(join)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &positions); err != nil {
		return nil, err
	}
	return positions, err
}

func (p *positionTestStore) ListInstrumentNamesByDate(ctx context.Context, time time.Time, status int64) (names []string, err error) {
	positions, err := p.loadTestPositions(time.Format(time2.RFC3339Date))
	if err != nil {
		log.Println(err)
	}
	var (
		nameMap = make(map[string]int, 0)
	)
	date := core.InstrumentDate(time)
	// eg. 1AUG21, rather than 11AUG21, should add '-' prefix
	if len(date) == 6 {
		date = "-" + date
	}
	if status == -1 {
		for _, e := range positions {
			if strings.Contains(e.InstrumentName, date) {
				if _, ok := nameMap[e.InstrumentName]; !ok {
					nameMap[e.InstrumentName]++
				}
			}
		}
	} else {
		for _, e := range positions {
			if strings.Contains(e.InstrumentName, date) && e.Status == int8(status) {
				if _, ok := nameMap[e.InstrumentName]; !ok {
					nameMap[e.InstrumentName]++
				}
			}
		}
	}

	for k := range nameMap {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}

func (p *positionTestStore) ListNamesByDateAndDeliveryType(ctx context.Context, time time.Time, status int64, deliveryType string) (names []string, err error) {
	allNames, err := p.ListInstrumentNamesByDate(ctx, time, status)
	if err != nil {
		return names, err
	}
	for _, e := range allNames {

		instrument, err := core.ParseInstrument(e)
		if err != nil {
			return names, err
		}
		if instrument.DeliveryType == deliveryType {
			names = append(names, e)
		}
	}
	return names, nil
}

func (p *positionTestStore) ListExercisablePosition(ctx context.Context, instrumentName string) (positions []*core.Position, err error) {
	return nil, nil
}

func (p *positionTestStore) ListMonthlyPremium(ctx context.Context, uid string, year int) (monthlyProfit []core.MonthlyProfit, err error) {
	return monthlyProfit, nil
}

func (p *positionTestStore) ListMonthlyUnderlying(ctx context.Context, uid string, year int) (monthlyProfit []core.MonthlyProfit, err error) {
	return monthlyProfit, nil
}

func (p *positionTestStore) ListSidePositions(ctx context.Context, instrumentName, side string, statusSet []int64) (positions []*core.Position, err error) {
	instrument, _ := core.ParseInstrument(instrumentName)
	testPositions, err := p.loadTestPositions(instrument.ExpirationDate.Format(time2.RFC3339Date))
	if err != nil {
		return nil, err
	}
	var (
		descPositions = make([]*core.Position, 0)
	)
	for _, e := range testPositions {
		if e.InstrumentName == instrumentName && e.Side == side {
			var statusPass = false
			for _, status := range statusSet {
				if e.Status == int8(status) {
					statusPass = true
					break
				}
			}
			if statusPass {
				descPositions = append(descPositions, e)
			}
		}
	}
	return descPositions, nil
}

func (p *positionTestStore) ListPositionTrades(uid, instrumentName, side string) (tradeList []*core.Trade, err error) {
	return nil, nil
}

func (p *positionTestStore) FindPositionStatusById(id string) (status int8, err error) {
	return 0, nil
}

func (p *positionTestStore) FindPositionByUidAndName(uid, instrumentName string) (position *core.Position, err error) {
	return nil, nil
}

func (p *positionTestStore) ListUserPositions(current, size int64, uid, status, order string) (list []*core.Position, total, pages int64, err error) {
	return nil, 0, 0, nil
}

func (p *positionTestStore) UpdateExercisedSize(id int64, size float64, updatedAt time.Time) error {
	return nil
}

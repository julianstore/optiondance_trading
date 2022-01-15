package position

import (
	"context"
	"fmt"
	"github.com/MixinNetwork/go-number"
	"gorm.io/gorm"
	"option-dance/core"
	"option-dance/pkg/util"
	"option-dance/store"
	"strconv"
	"time"
)

func New(db *gorm.DB) core.PositionStore {
	return &positionStore{db: db}
}

type positionStore struct {
	db *gorm.DB
}

func (s *positionStore) FindByPositionId(ctx context.Context, positionId string) (position *core.Position, err error) {
	if err = s.db.Model(core.Position{}).Where("position_id=?", positionId).Find(&position).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return position, nil
}

func (s *positionStore) ExercisePositionWithID(ctx context.Context, positionId, instrumentName string, updateTime time.Time) (err error) {
	if err = s.db.Model(core.Position{}).Where("position_id=?", positionId).Updates(map[string]interface{}{
		"status":     core.PositionStatusExerciseRequested,
		"updated_at": updateTime,
	}).Error; store.CheckErr(err) != nil {
		return err
	}
	//update asks positions
	if err = s.db.Model(core.Position{}).Where("instrument_name=? and side = 'ASK' ", instrumentName).Updates(map[string]interface{}{
		"status":     core.PositionStatusExerciseRequested,
		"updated_at": updateTime,
	}).Error; store.CheckErr(err) != nil {
		return err
	}
	return nil
}

func (s *positionStore) AutoExercisePositionWithName(ctx context.Context, instrumentName string, updateTime time.Time) (err error) {
	//update asks and bid positions
	if err = s.db.Model(core.Position{}).WithContext(ctx).Where("instrument_name=? ", instrumentName).Updates(map[string]interface{}{
		"status":         core.PositionStatusExerciseRequested,
		"exercised_size": gorm.Expr("abs(size)"),
		"updated_at":     updateTime,
	}).Error; store.CheckErr(err) != nil {
		return err
	}
	return nil
}

func (s *positionStore) ListExercisablePosition(ctx context.Context, instrumentName string) (positions []*core.Position, err error) {
	if err := s.db.Model(core.Position{}).
		Where("instrument_name = ? and side = 'BID' ", instrumentName).
		Find(&positions).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return positions, nil
}

// FindInstrumentNamesByDate (status = -1 || status < 0) means do not filter the position status
func (s *positionStore) ListInstrumentNamesByDate(ctx context.Context, time time.Time, status int64) (names []string, err error) {
	date := core.InstrumentDate(time)
	// eg. 1AUG21, rather than 11AUG21, should add '-' prefix
	if len(date) == 6 {
		date = "-" + date
	}
	if status < 0 {
		queryString := "SELECT instrument_name FROM `position` where settlement = 0 AND instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).
			Raw(queryString, "%"+date+"%").
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	} else {
		queryString := "SELECT instrument_name FROM `position` where settlement = 0 AND status = ? and instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).
			Raw(queryString, status, "%"+date+"%").
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	}
	return names, nil
}

// FindInstrumentNamesByDate (status = -1 || status < 0) means do not filter the position status
func (s *positionStore) ListNamesByDateAndDeliveryType(ctx context.Context, time time.Time, status int64, deliveryType string) (names []string, err error) {
	date := core.InstrumentDate(time)
	// eg. 1AUG21, rather than 11AUG21, should add '-' prefix
	if len(date) == 6 {
		date = "-" + date
	}
	if status < 0 {
		queryString := "SELECT instrument_name FROM `position` where delivery_type = ? and settlement = 0 AND instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).
			Raw(queryString, deliveryType, "%"+date+"%").WithContext(ctx).
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	} else {
		queryString := "SELECT instrument_name FROM `position` where delivery_type = ? and settlement = 0 AND status = ? and instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).WithContext(ctx).
			Raw(queryString, deliveryType, status, "%"+date+"%").
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	}
	return names, nil
}

// FindInstrumentNamesByDate (status = -1 || status < 0) means do not filter the position status
func (s *positionStore) ListNamesByDateAndSide(ctx context.Context, time time.Time, status int64, side, deliveryType string) (names []string, err error) {
	date := core.InstrumentDate(time)
	// eg. 1AUG21, rather than 11AUG21, should add '-' prefix
	if len(date) == 6 {
		date = "-" + date
	}
	if status < 0 {
		queryString := "SELECT instrument_name FROM `position` where side = ? and delivery_type = ? AND settlement = 0 AND instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).WithContext(ctx).
			Raw(queryString, side, deliveryType, "%"+date+"%").
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	} else {
		queryString := "SELECT instrument_name FROM `position` where side = ? and delivery_type = ? AND settlement = 0 AND status = ? and instrument_name like ? group by instrument_name"
		if err = s.db.Model(core.Position{}).WithContext(ctx).
			Raw(queryString, side, deliveryType, status, "%"+date+"%").
			Scan(&names).Error; store.CheckErr(err) != nil {
			return nil, err
		}
	}
	return names, nil
}

func (s *positionStore) ListMonthlyPremium(ctx context.Context, uid string, year int) (monthlyProfit []core.MonthlyProfit, err error) {
	querySql := `SELECT date_format(p.expiration_date,'%Y%m')as month,
					sum(p.funds) as amount
					FROM position p
					where p.user_id = ?
					and p.side = 'ASK'
					group by month`
	if err = s.db.Model(core.Position{}).
		Raw(querySql, uid).Find(&monthlyProfit).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return
}

func (s *positionStore) ListMonthlyUnderlying(ctx context.Context, uid string, year int) (monthlyProfit []core.MonthlyProfit, err error) {
	querySql := `SELECT date_format(p.expiration_date,'%Y%m')as month,
					sum(p.exercised_size) as amount
					FROM position p
					where p.user_id = ?
					and p.side = 'ASK'
					group by month`
	if err = s.db.Model(core.Position{}).
		Raw(querySql, uid).Scan(&monthlyProfit).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return
}

func (s *positionStore) ListSidePositions(ctx context.Context, instrumentName, side string, statusSet []int64) (positions []*core.Position, err error) {
	if err = s.db.Model(core.Position{}).WithContext(ctx).
		Where("instrument_name = ? AND side = ? AND status in (?)", instrumentName, side, statusSet).Find(&positions).Error; store.CheckErr(err) != nil {
		return nil, err
	}
	return positions, nil
}

func (s *positionStore) FindByPositionIdAndUid(ctx context.Context, uid, positionId string) (position *core.Position, err error) {
	if err = s.db.Model(core.Position{}).WithContext(ctx).Where("position_id=?", positionId).Find(&position).Error; err != nil {
		return nil, err
	}
	position.PositionNumString = strconv.Itoa(int(position.PositionNum))

	tradeList, err := s.ListPositionTrades(uid, position.InstrumentName, "")
	if err != nil {
		return nil, err
	}
	position.TradeList = tradeList
	funds, err := s.positionInitialFunds(tradeList, position.Side)
	if err != nil {
		return nil, err
	}
	position.InitialFunds = funds

	return position, nil
}

func (s *positionStore) positionInitialFunds(tradeList []*core.Trade, side string) (funds float64, err error) {
	var f number.Decimal
	for _, e := range tradeList {
		if e.Side == side {
			f = f.Add(number.FromString(e.Price).Mul(number.FromString(e.Amount)))
		}
	}
	return f.Float64(), nil
}

func (s *positionStore) ListPositionTrades(uid, instrumentName, side string) (tradeList []*core.Trade, err error) {
	qs := fmt.Sprintf("user_id='%s' AND instrument_name = '%s' AND side = '%s'", uid, instrumentName, side)
	if len(side) == 0 {
		qs = fmt.Sprintf("user_id='%s' AND instrument_name = '%s'", uid, instrumentName)
	}
	if err = s.db.Model(core.Trade{}).Where(qs).Find(&tradeList).Error; err != nil {
		return nil, err
	}
	return tradeList, nil
}

func (s *positionStore) FindPositionStatusById(id string) (status int8, err error) {
	var position *core.Position
	if err = s.db.Model(core.Position{}).Where("position_id=?", id).Select("status", "id").Find(&position).Error; err != nil {
		return core.PositionStatusNotExercised, err
	}
	if position.ID > 0 {
		return position.Status, nil
	}
	return core.PositionStatusNotExercised, nil
}

func (s *positionStore) FindPositionByUidAndName(uid, instrumentName string) (position *core.Position, err error) {
	if err = s.db.Model(core.Position{}).Where("user_id = ? AND instrument_name=?", uid, instrumentName).Find(&position).Error; err != nil {
		return nil, err
	}
	position.PositionNumString = strconv.Itoa(int(position.PositionNum))
	return position, nil
}

func (s *positionStore) ListUserPositions(current, size int64, uid, status, order string) (list []*core.Position, total, pages int64, err error) {
	orderCondition := "expiration_date"
	if order == "desc" || order == "asc" {
		orderCondition += " " + order
	}
	orderCondition += ", strike_price"
	queryString := fmt.Sprintf("`user_id` = '%s' AND (`status` = '%d' OR settlement = 1) ", uid, core.PositionStatusExercised)
	if status == "open" {
		queryString = fmt.Sprintf("`user_id` = '%s' AND settlement = 0 AND (`status` = '%d' OR `status` = '%d')",
			uid, core.PositionStatusNotExercised, core.PositionStatusExerciseRequested)
	}
	if err = s.db.Model(core.Position{}).Where(queryString).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}
	pages = util.GetPages(total, size)
	if err = s.db.Model(core.Position{}).Where(queryString).
		Offset(int((current - 1) * size)).Limit(int(size)).Order(orderCondition).Find(&list).Error; err != nil {
		return nil, 0, 0, err
	}
	for _, e := range list {
		e.PositionNumString = strconv.Itoa(int(e.PositionNum))
		trades, err := s.ListPositionTrades(uid, e.InstrumentName, e.Side)
		if err != nil {
			return nil, 0, 0, err
		}
		funds, err := s.positionInitialFunds(trades, e.Side)
		if err != nil {
			return nil, 0, 0, err
		}
		e.InitialFunds = funds
	}

	return
}

func (s *positionStore) UpdateExercisedSize(id int64, size float64, updatedAt time.Time) error {
	//update exercised size
	if err := s.db.Model(core.Position{}).
		Where("id = ?", id).Updates(map[string]interface{}{
		"exercised_size": size,
		"updated_at":     updatedAt,
	}).Error; store.CheckErr(err) != nil {
		return err
	}
	return nil
}

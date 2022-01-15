package market

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"option-dance/core"
	"time"
)

func NewRedisStore(c *redis.Client) core.MarketStore {
	return &redisStore{client: c}
}

type redisStore struct {
	client *redis.Client
}

const (
	PipelineLimit = 200
)

//KEY format: BTC-18JUN21-42000-P
func (r *redisStore) Create(market core.OptionMarket) error {
	dto, err := core.ConvertToMarketDTO(market)
	if err != nil {
		return err
	}
	err = r.Save(dto)
	if err != nil {
		return err
	}
	return nil
}

func (r *redisStore) Save(dto core.OptionMarketDTO) error {
	marshal, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	r.client.Set(dto.InstrumentName, string(marshal), 24*time.Hour)
	return nil
}

func (r *redisStore) FindByInstrumentName(instrumentName string) (o core.OptionMarketDTO, err error) {
	s, _ := r.client.Get(instrumentName).Result()
	var dto core.OptionMarketDTO
	err = json.Unmarshal([]byte(s), &dto)
	if err != nil {
		return core.OptionMarketDTO{}, err
	}
	return
}

func (r *redisStore) ListByStrikePrice(price, side, optionType, deliveryType, quoteCurrency, baseCurrency string) (list []core.OptionMarketDTO, err error) {
	match := fmt.Sprintf("*%s*", price)
	result, _, err := r.client.Scan(0, match, PipelineLimit).Result()
	return r.pipelineGetNames(result)
}

func (r *redisStore) pipelineGetNames(nameList []string) (list []core.OptionMarketDTO, err error) {
	pipeline := r.client.Pipeline()
	for _, e := range nameList {
		pipeline.Get(e)
	}
	exec, err := pipeline.Exec()
	if err != nil {
		return nil, nil
	}
	for _, e := range exec {
		var o core.OptionMarketDTO
		s, _ := e.(*redis.StringCmd).Result()
		err := json.Unmarshal([]byte(s), &o)
		if err != nil {
			return list, nil
		}
		list = append(list, o)
	}
	return list, nil
}

func (r *redisStore) ListMarketStrikePrice(side, optionType, deliveryType, quoteCurrency, baseCurrency string) (optionMarkets []core.OptionMarketDTO, err error) {
	var wantSide = "ASK"
	if side == "ASK" {
		wantSide = "BID"
	}
	members, err := r.client.SMembers(wantSide).Result()
	if err != nil {
		return nil, err
	}
	return r.pipelineGetNames(members)
}

func (r *redisStore) UpdateBidAsks(instrumentName string, bids, asks string) (rowsAffected int64, err error) {
	if bids == "[]" {
		r.client.SRem("BID", instrumentName)
	} else {
		r.client.SAdd("BID", instrumentName)
	}

	if asks == "[]" {
		r.client.SRem("ASK", instrumentName)
	} else {
		r.client.SAdd("ASK", instrumentName)
	}

	i, err := core.ParseInstrument(instrumentName)
	if err != nil {
		return 1, err
	}
	market := core.OptionMarket{
		InstrumentName:      instrumentName,
		StrikePrice:         i.StrikePrice,
		ExpirationDate:      i.ExpirationDate,
		ExpirationDateStr:   i.ExpirationDate.Format("2006-01-02"),
		ExpirationTimestamp: i.ExpirationTimestamp,
		OptionType:          i.OptionType,
		BaseCurrency:        i.BaseCurrency,
		Bids:                bids,
		Asks:                asks,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Status:              0,
	}
	err = r.Create(market)
	if err != nil {
		return 1, err
	}
	return 1, nil
}

func (r *redisStore) DeleteByInstrumentName(ctx context.Context, instrumentName string) error {
	return nil
}

func (r *redisStore) ListDates(ctx context.Context, deliveryType, quoteCurrency, baseCurrency string) (dateList []string, err error) {
	return nil, err
}

func (r *redisStore) ListByDate(ctx context.Context, date string) (dateList []core.FullMarketByDateDTO, err error) {
	return nil, err
}

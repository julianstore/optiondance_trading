package position

import (
	"context"
	"option-dance/cmd/config"
	"option-dance/core"
	time2 "option-dance/pkg/time"
	"testing"
	"time"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	config.InitTestConfig("../../config/od_dev_config/node1.yaml")
	m.Run()
}

func getPositionStore() core.PositionStore {
	store := New(config.Db)
	return store
}

func TestPositionTestStore_ListInstrumentNamesByDate(t *testing.T) {
	config.InitTestConfig("../../config/od_dev_config/node1.yaml")
	positionTestStore := NewPositionTestStore()
	ctx := context.Background()
	parse, err := time.Parse(time2.RFC3339Date, "2021-09-10")
	if err != nil {
		t.Error(err)
	}
	date, err := positionTestStore.ListInstrumentNamesByDate(ctx, parse, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(date)
}

func TestPositionStore_ExercisePositionWithName(t *testing.T) {
	store := getPositionStore()
	err := store.AutoExercisePositionWithName(ctx, "P-pUSD-BTC-6SEP21-40000-P", time.Now())
	if err != nil {
		t.Error(err)
	}
}

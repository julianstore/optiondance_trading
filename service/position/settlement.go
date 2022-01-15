package position

import (
	"context"
	"fmt"
	"option-dance/cmd/config"
	"option-dance/core"
	time2 "option-dance/pkg/time"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//Query the instrument on the day, and then calculate the number of right positions, the number of right positions,
//the proportion, and the number of obligations of each instrument.

//Generating transfers: assigning rules and distributing proportionally to the users of the obligatory warehouse transfers (the target and the unexercised part (proportionally assigned)
//and the remaining margin), and transfers the exercising part of the deposit to the users of the right warehouse
func (p *positionService) Settlement(ctx context.Context, asset, dateString, deliveryType string, settledAt time.Time) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	zap.L().Info("settlement:start")

	result, err := p.SettlementGather(ctx, asset, dateString, deliveryType, true, settledAt)
	if err != nil {
		zap.L().Info("settlement gather error", zap.Error(err))
		return err
	}

	go p.LogGatherResult(result)
	err = p.SettlementCommit(ctx, result.Transfers, result.ExercisedNameList, result.NameList, settledAt)
	if err != nil {
		zap.L().Info("settlement commit error", zap.Error(err))
	}

	zap.L().Info("settlement done")
	return nil
}

func (p *positionService) SettlementGather(ctx context.Context, asset, dateString, deliveryType string, updateDb bool, settledAt time.Time) (core.GatherResult, error) {

	var (
		nameList          []string
		exercisedNameList []string
		transfers         []*core.Transfer
	)

	date, err := time.Parse(time2.RFC3339Date, dateString)
	if err != nil {
		return core.GatherResult{}, err
	}

	//get settlement nameList
	nameList, exercisedNameList, err = p.GetSettlementNameList(ctx, date, deliveryType)
	if err != nil {
		zap.L().Error("GetSettlementNameList error", zap.Error(err))
		return core.GatherResult{}, err
	}

	//generate transfer
	transfers, err = p.GenSettlementTransfer(ctx, deliveryType, nameList, updateDb, settledAt)
	if err != nil {
		zap.L().Error("GenSettlementTransfer error", zap.Error(err))
		return core.GatherResult{}, err
	}

	return core.GatherResult{
		NameList:          nameList,
		ExercisedNameList: exercisedNameList,
		Transfers:         transfers,
	}, nil
}

func (p *positionService) GetSettlementNameList(ctx context.Context, date time.Time, deliveryType string) (nameList, exercisedNameList []string, err error) {

	nameList, err = p.positionStore.ListNamesByDateAndDeliveryType(ctx, date, core.PositionStatusALL, deliveryType)
	if err != nil {
		return nil, nil, err
	}

	exercisedNameList, err = p.positionStore.ListNamesByDateAndSide(ctx, date, core.PositionStatusExerciseRequested, core.PositionSideBid, deliveryType)
	if err != nil {
		return nil, nil, err
	}

	zap.L().Info("Cron-Settlement", zap.Strings("nameList", nameList), zap.Strings("exercisedNameList", exercisedNameList))
	return
}

func (p *positionService) GenSettlementTransfer(ctx context.Context, deliveryType string, nameList []string, updateDb bool, settledAt time.Time) (transferList []*core.Transfer, err error) {

	var (
		callCashDeliveryOptionNameList     = make([]string, 0)
		callPhysicalDeliveryOptionNameList = make([]string, 0)
		putCashDeliveryOptionNameList      = make([]string, 0)
		putPhysicalDeliveryOptionNameList  = make([]string, 0)
		dateString                         = ""
	)
	for i, e := range nameList {
		info, err := core.ParseInstrument(e)
		if err != nil {
			return nil, err
		}

		if deliveryType == core.DeliveryTypeCash {
			if info.OptionType == core.OptionTypePUT && info.DeliveryType == core.DeliveryTypeCash {
				putCashDeliveryOptionNameList = append(putCashDeliveryOptionNameList, e)
			}
			if info.OptionType == core.OptionTypeCALL && info.DeliveryType == core.DeliveryTypeCash {
				callCashDeliveryOptionNameList = append(callCashDeliveryOptionNameList, e)
			}
		}
		if deliveryType == core.DeliveryTypePhysical {
			if info.OptionType == core.OptionTypePUT && info.DeliveryType == core.DeliveryTypePhysical {
				putPhysicalDeliveryOptionNameList = append(putPhysicalDeliveryOptionNameList, e)
			}
			if info.OptionType == core.OptionTypeCALL && info.DeliveryType == core.DeliveryTypePhysical {
				callPhysicalDeliveryOptionNameList = append(callPhysicalDeliveryOptionNameList, e)
			}
		}

		if i == 0 {
			dateString = info.ExpirationDate.Format(time2.RFC3339Date)
		}
	}

	sort.Strings(callCashDeliveryOptionNameList)
	sort.Strings(callPhysicalDeliveryOptionNameList)
	sort.Strings(putCashDeliveryOptionNameList)
	sort.Strings(putPhysicalDeliveryOptionNameList)

	var deliveryPrice number.Decimal
	if len(callCashDeliveryOptionNameList) > 0 || len(putCashDeliveryOptionNameList) > 0 {
		price, err := p.deliveryPriceStore.ReadPrice(ctx, core.BTC, dateString)
		if err != nil {
			return nil, err
		}
		deliveryPrice = number.FromString(price)
		zap.L().Info("delivery price", zap.String("price", deliveryPrice.String()))
		if deliveryPrice.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("delivery price not correct: %s", deliveryPrice.String())
		}
	}

	if transferList, err = p.appendTransfers(ctx, putCashDeliveryOptionNameList, core.DeliveryTypeCash, core.OptionTypePUT,
		deliveryPrice, updateDb, true, settledAt, transferList); err != nil {
		return nil, err
	}
	if transferList, err = p.appendTransfers(ctx, callCashDeliveryOptionNameList, core.DeliveryTypeCash, core.OptionTypeCALL,
		deliveryPrice, updateDb, true, settledAt, transferList); err != nil {
		return nil, err
	}
	if transferList, err = p.appendTransfers(ctx, putPhysicalDeliveryOptionNameList, core.DeliveryTypePhysical, core.OptionTypePUT,
		deliveryPrice, updateDb, false, settledAt, transferList); err != nil {
		return nil, err
	}
	if transferList, err = p.appendTransfers(ctx, callPhysicalDeliveryOptionNameList, core.DeliveryTypePhysical, core.OptionTypeCALL,
		deliveryPrice, updateDb, false, settledAt, transferList); err != nil {
		return nil, err
	}
	return transferList, nil
}

func (p *positionService) appendTransfers(ctx context.Context, nameList []string, deliveryType, optionType string, deliveryPrice number.Decimal,
	updateDb, autoExercise bool, settledAt time.Time, transferList []*core.Transfer) ([]*core.Transfer, error) {

	for _, name := range nameList {
		var askPositions []*core.Position
		instrument, strikePrice, exerciseBidPositions, exerciseRatio, err := p.prepareGenTransfers(ctx, name, deliveryPrice, settledAt, autoExercise)
		if err != nil {
			return nil, err
		}

		if deliveryType == core.DeliveryTypeCash && optionType == core.OptionTypePUT {
			//bid exercised user transfers
			bidTransfers, err := p.CashBidPutPositionTransfers(exerciseBidPositions, updateDb, strikePrice, deliveryPrice, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, bidTransfers...)

			//ask exercised user transfer
			askPositions, err = p.positionStore.ListSidePositions(ctx, name, "ASK", []int64{core.PositionStatusNotExercised, core.PositionStatusExerciseRequested})
			if err != nil {
				return nil, err
			}
			askTransfers, err := p.CashAskPutPositionTransfers(askPositions, updateDb, strikePrice, deliveryPrice, exerciseRatio, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, askTransfers...)
		} else if deliveryType == core.DeliveryTypeCash && optionType == core.OptionTypeCALL {
			//bid exercised user transfers
			bidTransfers, err := p.CashBidCallPositionTransfers(exerciseBidPositions, updateDb, strikePrice, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, bidTransfers...)

			//ask exercised user transfer
			askPositions, err = p.positionStore.ListSidePositions(ctx, name, "ASK", []int64{core.PositionStatusNotExercised, core.PositionStatusExerciseRequested})
			if err != nil {
				return nil, err
			}
			askTransfers, err := p.CashAskCallPositionTransfers(askPositions, updateDb, strikePrice, exerciseRatio, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, askTransfers...)
		} else if deliveryType == core.DeliveryTypePhysical && optionType == core.OptionTypePUT {
			//bid exercised user transfers
			bidTransfers, err := p.PhysicalBidPutPositionTransfers(exerciseBidPositions, updateDb, strikePrice, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, bidTransfers...)

			//ask exercised user transfer
			askPositions, err = p.positionStore.ListSidePositions(ctx, name, "ASK", []int64{core.PositionStatusNotExercised, core.PositionStatusExerciseRequested})
			if err != nil {
				return nil, err
			}
			askTransfers, err := p.PhysicalAskPutPositionTransfers(askPositions, updateDb, strikePrice, exerciseRatio, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, askTransfers...)
		} else if deliveryType == core.DeliveryTypePhysical && optionType == core.OptionTypeCALL {
			bidTransfers, err := p.PhysicalBidCallPositionTransfers(exerciseBidPositions, updateDb, strikePrice, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, bidTransfers...)

			//ask exercised user transfer
			askPositions, err = p.positionStore.ListSidePositions(ctx, name, "ASK", []int64{core.PositionStatusNotExercised, core.PositionStatusExerciseRequested})
			if err != nil {
				return nil, err
			}
			askTransfers, err := p.PhysicalAskCallPositionTransfers(askPositions, updateDb, strikePrice, exerciseRatio, instrument, settledAt)
			if err != nil {
				return nil, err
			}
			transferList = append(transferList, askTransfers...)
		}
	}
	return transferList, nil
}

func (p *positionService) prepareGenTransfers(ctx context.Context, name string, deliveryPrice number.Decimal, settledAt time.Time, autoExercise bool) (instrument core.OptionInfo, strikePrice number.Decimal,
	exerciseBidPositions []*core.Position, exerciseRatio string, err error) {

	instrument, err = core.ParseInstrument(name)
	if err != nil {
		return instrument, number.Zero(), nil, "0", err
	}
	strikePrice = number.FromString(strconv.Itoa(int(instrument.StrikePrice)))

	//calculate exercise ratio
	exerciseBidPositions, exerciseRatio, err = p.CalcExerciseRatio(ctx, name, instrument.DeliveryType, deliveryPrice, strikePrice, settledAt, autoExercise)
	if err != nil {
		return instrument, number.Zero(), nil, "0", err
	}
	return
}

func (p *positionService) SettlementCommit(ctx context.Context, transferList []*core.Transfer, exercisedNameList, nameList []string, settledAt time.Time) error {

	//batch insert all transfer & update all position status
	if err := config.Db.Transaction(func(tx *gorm.DB) error {
		if len(transferList) > 0 {
			if err := tx.Model(core.Transfer{}).Create(&transferList).Error; err != nil {
				return err
			}
		}
		for _, e := range nameList {
			tx.Model(core.Position{}).Where("instrument_name = ? AND status = ?", e, core.PositionStatusExerciseRequested).
				Updates(map[string]interface{}{
					"status":     core.PositionStatusExercised,
					"updated_at": settledAt,
				})
		}
		tx.Model(core.Position{}).Where("instrument_name in (?)", nameList).
			Updates(map[string]interface{}{
				"settlement": 1,
				"updated_at": settledAt,
			})
		return nil
	}); err != nil {
		return err
	}
	zap.L().Info("settlement commit success", zap.Int("transfersLen", len(transferList)))
	return nil
}

func (p *positionService) LogGatherResult(result core.GatherResult) {
	zap.L().Info("Instruments", zap.Strings("names", result.NameList), zap.Strings("exercisedNames", result.ExercisedNameList))
	for _, e := range result.Transfers {
		user, _ := p.userStore.FindByMixinId(e.UserId)
		zap.L().Info("transfer", zap.String("user", e.UserId+"|"+user.Nickname), zap.String("amount", e.Amount), zap.String("source", e.Source))
	}
}

func ToAmountString(d decimal.Decimal) (amount string) {
	s := d.String()
	if strings.Contains(s, ".") {
		split := strings.Split(s, ".")
		if len(split) != 2 {
			amount = s
		} else {
			if len(split[1]) <= 8 {
				amount = s
			} else {
				amount = split[0] + "." + split[1][:8]
			}
		}
	} else {
		amount = s
	}
	return number.FromString(amount).String()
}

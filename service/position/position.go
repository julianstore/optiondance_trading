package position

import (
	"context"
	"option-dance/core"
	time2 "option-dance/pkg/time"
	"option-dance/pkg/util"
	"strconv"
	"sync"
	"time"

	"github.com/MixinNetwork/go-number"
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(positionStore core.PositionStore,
	tradeStore core.TradeStore,
	transferService core.TransferService,
	userStore core.UserStore,
	deliveryPriceStore core.DeliveryPriceStore,
) core.PositionService {
	return &positionService{
		positionStore:      positionStore,
		tradeStore:         tradeStore,
		mutex:              sync.Mutex{},
		transferService:    transferService,
		userStore:          userStore,
		deliveryPriceStore: deliveryPriceStore,
	}
}

func NewTest(positionStore core.PositionStore, deliveryPriceStore core.DeliveryPriceStore) core.PositionService {
	return &positionService{
		positionStore:      positionStore,
		tradeStore:         nil,
		mutex:              sync.Mutex{},
		transferService:    nil,
		userStore:          nil,
		deliveryPriceStore: deliveryPriceStore,
	}
}

type positionService struct {
	positionStore      core.PositionStore
	tradeStore         core.TradeStore
	transferService    core.TransferService
	userStore          core.UserStore
	mutex              sync.Mutex
	deliveryPriceStore core.DeliveryPriceStore
}

func (p *positionService) PositionDetail(ctx context.Context, uid, positionId string) (dto *core.PositionDTO, err error) {
	position, err := p.positionStore.FindByPositionIdAndUid(ctx, uid, positionId)
	if err != nil {
		return nil, err
	}
	info, err := p.SettlementInfo(ctx, position)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	expiryTs := position.ExpirationTimestamp
	expired := core.ExpiredPosition(now, expiryTs)
	exercisable := core.Exercisable(now, expiryTs)
	var btnStatus = ""
	if govalidator.IsIn(position.Side, core.PositionSideAsk, core.PositionSideZero) || expired ||
		position.Settlement == 1 || position.DeliveryType == core.DeliveryTypeCash {
		btnStatus = core.ExerciseBtnStatusHidden
	} else if !exercisable {
		btnStatus = core.ExerciseBtnStatusWaitExercise
	} else if exercisable && position.Status == core.PositionStatusNotExercised {
		btnStatus = core.ExerciseBtnStatusOpenExercise
	} else if exercisable && position.Status == core.PositionStatusExerciseRequested {
		btnStatus = core.ExerciseBtnStatusWaitSettlement
	}
	dto = &core.PositionDTO{
		Position:       position,
		SettlementInfo: info,
		Exercisable:    exercisable,
		Expiry:         expired,
		BtnStatus:      btnStatus,
	}
	return dto, nil
}

func (p *positionService) SettlementInfo(ctx context.Context, position *core.Position) (info core.SettlementInfo, err error) {
	format := position.ExpirationDate.Format(time2.RFC3339Date)
	price, err := p.deliveryPriceStore.ReadPrice(ctx, position.BaseCurrency, format)
	if err != nil {
		return info, err
	}
	if price == "" {
		return info, nil
	}
	size := decimal.NewFromFloat(position.Size).Abs()
	deliveryPrice := decimal.RequireFromString(price)
	strikePrice := decimal.RequireFromString(position.StrikePrice)
	refundMargin := size.Mul(strikePrice)
	bidEarnings := decimal.Zero
	exercised := false
	if strikePrice.GreaterThanOrEqual(deliveryPrice) {
		refundMargin = size.Mul(deliveryPrice)
		diffPrice := strikePrice.Sub(deliveryPrice)
		bidEarnings = size.Mul(diffPrice)
		exercised = true
	}
	if strikePrice.LessThan(deliveryPrice) {
		size = decimal.Zero
	}

	info = core.SettlementInfo{
		UnderlyingPrice: deliveryPrice.String(),
		RefundMargin:    refundMargin.String(),
		Size:            size.String(),
		BidEarnings:     bidEarnings.String(),
		Exercised:       exercised,
	}
	return
}

func (p *positionService) UpdatePositionWithTrade(tx *gorm.DB, trade *core.Trade) (err error) {
	var position *core.Position
	if err := tx.Model(core.Position{}).
		Where("user_id = ? AND instrument_name = ?", trade.UserID, trade.InstrumentName).
		Find(&position).Error; err != nil {
		return nil
	}
	var newSize, tradeFunds, newFunds, tradeMargin, newMargin number.Decimal
	zero := number.Zero()
	instrument, err := core.ParseInstrument(trade.InstrumentName)
	if err != nil {
		return err
	}
	if instrument.OptionType == core.OptionTypePUT {
		tradeMargin = number.FromString(trade.Amount).Mul(number.FromFloat(float64(instrument.StrikePrice)))
	} else {
		tradeMargin = number.FromString(trade.Amount)
	}
	//update exists position
	tradeFunds = number.FromString(trade.Amount).Mul(number.FromString(trade.Price))
	if position != nil && position.ID > 0 {
		//increase positions
		if trade.Side == core.PositionSideBid {
			newSize = number.FromFloat(position.Size).Add(number.FromString(trade.Amount))
			newFunds = number.FromFloat(position.Funds).Sub(tradeFunds)
			oldMargin := number.FromFloat(position.Margin)
			newMargin = oldMargin
			//You need to hold a good buy or a steady win before you can use the guaranteed sell, and the guaranteed buy to close the position
			if number.FromFloat(position.Size).Cmp(zero) < 0 {
				refundMargin := tradeMargin
				if oldMargin.Cmp(zero) <= 0 {
					refundMargin = zero
				} else if oldMargin.Cmp(refundMargin) <= 0 {
					refundMargin = oldMargin
				}
				newMargin = oldMargin.Sub(refundMargin)
				//refund margin
				if trade.Status == core.TradeStatusOpen {
					if err = tx.Model(core.Trade{}).Where("instrument_name = ? AND user_id = ? AND side = ?", trade.InstrumentName, trade.UserID, core.PositionSideBid).
						Updates(map[string]interface{}{
							"updated_at": trade.CreatedAt,
							"status":     core.TradeStatusClosed,
						}).Error; err != nil {
						return err
					}
					err := p.transferService.CreateClosePositionTransfer(tx, refundMargin.Float64(), position.PositionID, position.UserID, instrument.OptionType, trade)
					if err != nil {
						return err
					}
				}
			}
		} else {
			//decrease positions
			newSize = number.FromFloat(position.Size).Sub(number.FromString(trade.Amount))
			newFunds = number.FromFloat(position.Funds).Add(tradeFunds)
			oldMargin := number.FromFloat(position.Margin)
			newMargin = oldMargin.Add(tradeMargin)
			//You need to hold the position to guarantee the sale and guarantee the buy, before you can use the excellent buy or the steady win to close the position
			var refundMargin number.Decimal
			if number.FromFloat(position.Size).Cmp(zero) > 0 {
				if newSize.Cmp(zero) >= 0 { //Excellent buy<=guaranteed sale, refund and excellent buy part
					refundMargin = tradeMargin
					newMargin = number.FromFloat(position.Margin)
				} else {
					refundMargin = number.FromFloat(position.Margin)
					newMargin = tradeMargin
				}
				if trade.Status == core.TradeStatusOpen {
					if err = tx.Model(core.Trade{}).Where("instrument_name = ? AND user_id = ? AND side = ?", trade.InstrumentName, trade.UserID, core.PositionSideAsk).
						Updates(map[string]interface{}{
							"updated_at": trade.CreatedAt,
							"status":     core.TradeStatusClosed,
						}).Error; err != nil {
						return err
					}
					err := p.transferService.CreateClosePositionTransfer(tx, refundMargin.Float64(), position.PositionID, position.UserID, instrument.OptionType, trade)
					if err != nil {
						return err
					}
				}
			}
		}
		position.Size = newSize.Float64()
		position.Funds = newFunds.Float64()
		position.Side = core.PositionSideFromSize(newSize)
		position.Margin = newMargin.Float64()
		if err = tx.Model(core.Position{}).Where("id = ?", position.ID).Updates(map[string]interface{}{
			"size":       position.Size,
			"funds":      position.Funds,
			"side":       position.Side,
			"status":     core.PositionStatusNotExercised,
			"margin":     position.Margin,
			"updated_at": trade.CreatedAt,
		}).Error; err != nil {
			return err
		}
	} else {
		//create a new position
		if trade.Side == core.PositionSideBid {
			newSize = zero.Add(number.FromString(trade.Amount))
			newFunds = zero.Sub(tradeFunds)
			newMargin = zero
		} else {
			newSize = zero.Sub(number.FromString(trade.Amount))
			newFunds = zero.Add(tradeFunds)
			newMargin = tradeMargin
		}
		position.Side = core.PositionSideFromSize(newSize)
		positionId := core.GetSettlementId(trade.UserID, trade.InstrumentName)
		pnum, err := util.GenSnowflakeIdInt64()
		if err != nil {
			return err
		}
		newPosition := core.Position{
			PositionID:          positionId,
			PositionNum:         pnum,
			UserID:              trade.UserID,
			InstrumentName:      trade.InstrumentName,
			Side:                position.Side,
			Type:                instrument.OptionType,
			DeliveryType:        instrument.DeliveryType,
			StrikePrice:         strconv.Itoa(int(instrument.StrikePrice)),
			OptionType:          instrument.OptionType,
			QuoteCurrency:       instrument.QuoteCurrency,
			BaseCurrency:        instrument.BaseCurrency,
			Size:                newSize.Float64(),
			Funds:               newFunds.Float64(),
			Margin:              newMargin.Float64(),
			Status:              core.PositionStatusNotExercised,
			ExpirationDate:      instrument.ExpirationDate,
			ExpirationTimestamp: instrument.ExpirationTimestamp,
			CreatedAt:           trade.CreatedAt,
			UpdatedAt:           trade.CreatedAt,
		}
		if err = tx.Model(core.Position{}).Create(&newPosition).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *positionService) AutoExercise(ctx context.Context, dateString string, deliveryPrice decimal.Decimal, exercisedAt time.Time) error {

	instrumentTime, err := time.Parse(time2.RFC3339Date, dateString)
	if err != nil {
		return err
	}
	names, err := p.positionStore.ListNamesByDateAndDeliveryType(ctx, instrumentTime, core.PositionStatusALL, core.DeliveryTypeCash)
	if err != nil {
		return err
	}
	for _, e := range names {
		instrument, err := core.ParseInstrument(e)
		if err != nil {
			return err
		}
		// https://www.notion.so/optiondance/15ea2bb17a8d4e4cb1d68f5da730cf69#78ce08e8d6ab4843ac38029c7db7f75e
		// if strike price greater than delivery price,should be auto exercised
		if instrument.DeliveryType == core.DeliveryTypeCash &&
			decimal.NewFromInt(instrument.StrikePrice).GreaterThan(deliveryPrice) {
			if err = p.positionStore.AutoExercisePositionWithName(ctx, e, exercisedAt); err != nil {
				return err
			}
			zap.L().Info("instrument auto exercised", zap.String("name", e))
		}
	}
	return nil
}

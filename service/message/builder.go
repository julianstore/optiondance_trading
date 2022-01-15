package message

import (
	"fmt"
	"option-dance/core"

	"go.uber.org/zap"
)

type TransferMessage struct {
	Source       string
	Title        string
	Instrument   string
	OptionType   string
	DeliveryType string
	Side         string
	Size         string
	Amount       string
	BaseCurrency string
}

type TransferMsgBuilder interface {
	New() TransferMsgBuilder
	SetSource(name string) TransferMsgBuilder
	SetTitle(title string) TransferMsgBuilder
	SetInstrument(side, instrument string) TransferMsgBuilder
	SetSide(side, optionType string) TransferMsgBuilder
	SetSize(side, totalSize, filledSize, remainingSize, baseAssetId, quoteAssetId string) TransferMsgBuilder
	SetAmount(side, amount, baseAssetId, quoteAssetId string) TransferMsgBuilder
	BuildMessage() string
}

type TransferMessageBuilder struct {
	msg *TransferMessage
}

func (b *TransferMessageBuilder) New() TransferMsgBuilder {
	b.msg = &TransferMessage{}
	return b
}

func (b *TransferMessageBuilder) SetSource(source string) TransferMsgBuilder {
	b.msg.Source = source
	return b
}

func (b *TransferMessageBuilder) SetTitle(title string) TransferMsgBuilder {
	b.msg.Title = fmt.Sprintf("【%s】", title)
	return b
}

func (b *TransferMessageBuilder) SetInstrument(side, instrument string) TransferMsgBuilder {
	i, err := core.ParseInstrument(instrument)
	if err != nil {
		zap.L().Error("ParseInstrument error", zap.String("name", instrument))
		return b
	}
	dateString := i.ExpirationDate.Format("2006年01月02日")
	typeCn := getTypeCn(side, i.OptionType)
	sideCn := getSideString(side, i.OptionType)
	b.msg.Instrument = fmt.Sprintf(`🎯 %s：%s %d %s %s %s`, typeCn, dateString, i.StrikePrice, core.GlobalQuoteCurrency, sideCn, i.BaseCurrency)
	b.msg.OptionType = i.OptionType
	b.msg.DeliveryType = i.DeliveryType
	b.msg.BaseCurrency = i.BaseCurrency
	return b
}

func getTypeCn(side, optionType string) (sideCn string) {
	if side == "ASK" && optionType == "PUT" {
		sideCn = "优买"
	}
	if side == "BID" && optionType == "PUT" {
		sideCn = "保卖"
	}
	if side == "ASK" && optionType == "CALL" {
		sideCn = "稳赢"
	}
	if side == "BID" && optionType == "CALL" {
		sideCn = "保买"
	}
	return
}

func getSideString(side, optionType string) (sideString string) {
	if side == "ASK" && optionType == "PUT" {
		sideString = "买入"
	}
	if side == "BID" && optionType == "PUT" {
		sideString = "卖出"
	}
	if side == "ASK" && optionType == "CALL" {
		sideString = "卖出"
	}
	if side == "BID" && optionType == "CALL" {
		sideString = "买入"
	}
	return
}

func (b *TransferMessageBuilder) SetSide(side, optionType string) TransferMsgBuilder {
	if side == "ASK" && optionType == "PUT" {
		b.msg.Side = fmt.Sprintf("- 📈 方向：%s", "优买")
	}
	if side == "BID" && optionType == "PUT" {
		b.msg.Side = fmt.Sprintf("- 📈 方向：%s", "保卖")
	}
	if side == "ASK" && optionType == "CALL" {
		b.msg.Side = fmt.Sprintf("- 📈 方向：%s", "稳赢")
	}
	if side == "BID" && optionType == "CALL" {
		b.msg.Side = fmt.Sprintf("- 📈 方向：%s", "保买")
	}
	return b
}

func (b *TransferMessageBuilder) SetSize(side, totalSize, filledSize, remainingSize, baseAssetId, quoteAssetId string) TransferMsgBuilder {
	var msg string = ""
	baseAsset := core.GetCurrencyByAsset(baseAssetId, true)
	switch b.msg.Source {
	case core.TransferSourceTradeConfirmed:
		msg = fmt.Sprintf("💎 数量: %s %s", filledSize, baseAsset)
		break
	case core.TransferSourceOrderCancelled, core.TransferSourceOrderInvalid:
		msg = fmt.Sprintf("💎 数量: %s %s", remainingSize, baseAsset)
		break
	case core.TransferSourcePositionExercised, core.TransferSourcePositionExercisedRefund, core.TransferSourcePositionExercisedInvalid:
		//sideString := getSideString(side, b.msg.OptionType)
		//if side == "BID" && b.msg.OptionType == "PUT" {
		//	msg = fmt.Sprintf("%s: %s %s", "💰 "+sideString, number.FromString(totalSize).Abs().String(), b.msg.BaseCurrency)
		//}
		break
	case core.TransferSourcePositionClosedRefund:
		//msg = fmt.Sprintf("Number of Closed Positions: %s %s", filledSize, baseAsset)
		break
	}
	b.msg.Size = msg
	return b
}

func (b *TransferMessageBuilder) SetAmount(side, amount, baseAssetId, quoteAssetId string) TransferMsgBuilder {
	var msgPrefix = "💰 "
	baseAsset, quoteAsset := core.GetCurrencyByAsset(baseAssetId, true), core.GetCurrencyByAsset(quoteAssetId, true)
	var asset string
	optionType := b.msg.OptionType
	switch b.msg.Source {
	case core.TransferSourceOrderCancelled, core.TransferSourceOrderInvalid:
		if side == "ASK" && optionType == "PUT" {
			msgPrefix += "解冻购币金: "
		}
		if side == "BID" && optionType == "PUT" {
			msgPrefix += "退回保卖费: "
		}
		if side == "ASK" && optionType == "CALL" {
			msgPrefix += "返还现金: "
		}
		if side == "BID" && optionType == "CALL" {
			msgPrefix += "返还现金: "
		}
		asset = quoteAsset
		break
	case core.TransferSourceTradeConfirmed:
		if side == "ASK" && optionType == "PUT" {
			msgPrefix += "获得现金: "
		}
		if side == "BID" && optionType == "PUT" {
			msgPrefix += "支付保卖费: "
		}
		if side == "ASK" && optionType == "CALL" {
			msgPrefix += "获得现金: "
		}
		if side == "BID" && optionType == "CALL" {
			msgPrefix += "支付保买费: "
		}
		asset = quoteAsset
		break
	case core.TransferSourcePositionExercised:
		if b.msg.DeliveryType == core.DeliveryTypePhysical {
			msgPrefix += "获得: "
			asset = baseAsset
		} else {
			if side == "ASK" {
				msgPrefix += "获得购币金: "
			} else {
				msgPrefix += "获得: "
			}
			asset = quoteAsset
		}
		break
	case core.TransferSourcePositionExercisedRefund:
		msgPrefix += "退回购币金: "
		asset = quoteAsset
		break
	case core.TransferSourcePositionExercisedInvalid:
		msgPrefix += "退还标的资产: "
		asset = baseAsset
		break
	case core.TransferSourcePositionClosedRefund:
		msgPrefix += "退回购币金: "
		asset = quoteAsset
		break
	}
	b.msg.Amount = fmt.Sprintf("%s %s %s", msgPrefix, amount, asset)
	return b
}

func (b *TransferMessageBuilder) BuildMessage() (msg string) {
	msg += b.msg.Title + "\n" + b.msg.Instrument + "\n"
	if len(b.msg.Size) > 0 {
		msg += b.msg.Size + "\n"
	}
	if len(b.msg.Amount) > 0 {
		msg += b.msg.Amount
	}

	//return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", b.msg.Title, b.msg.Instrument, b.msg.Side, b.msg.Size, b.msg.Amount)
	return
}

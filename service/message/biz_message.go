package message

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"io"
	"option-dance/cmd/config"
	"option-dance/core"
	"strings"
)

var (
	iconMaps = map[string]string{
		"CNB":  "https://mixin-images.zeromesh.net/0sQY63dDMkWTURkJVjowWY6Le4ICjAFuu3ANVyZA4uI3UdkbuOT5fjJUT82ArNYmZvVcxDXyNjxoOv0TAYbQTNKS=s128",
		"USDT": "https://mixin-images.zeromesh.net/ndNBEpObYs7450U08oAOMnSEPzN66SL8Mh-f2pPWBDeWaKbXTPUIdrZph7yj8Z93Rl8uZ16m7Qjz-E-9JFKSsJ-F=s128",
		"USDC": "https://mixin-images.zeromesh.net/w3Lb-pMrgcmmrzamf7FG_0c6Dkh3w_NRbysqzpuacwdVhMYSOtnX2zedWqiSG7JuZ3jd4xfhAJduQXY1rPidmywn=s128",
		"pUSD": "https://mixin-images.zeromesh.net/cH4GWuPXbzeZl6OOunpn7BxE25n3v8URwnNszs0FpZqv3OTlxP1zpzKw89VKTpBwWL-Ned1R36mmy1C4GMuPX1rL-PjfEJ2zby9V=s128",
		"BTC":  "https://mixin-images.zeromesh.net/HvYGJsV5TGeZ-X9Ek3FEQohQZ3fE9LBEBGcOcn4c4BNHovP4fW4YB97Dg5LcXoQ1hUjMEgjbl1DPlKg1TW7kK6XP=s128",
		"ETH":  "https://mixin-images.zeromesh.net/sdlLs6NIN7Qd_fM3z4k-s4GtamiHo4pafLve7P7hnHaLEPs2ld0FSscEzqjdytGY2q8e-AyfCNVfqXbXYYqb8cwUf-X3oI2jRTZN=s128",
		"XIN":  "https://mixin-images.zeromesh.net/UasWtBZO0TZyLTLCFQjvE_UYekjC7eHCuT_9_52ZpzmCC-X-NPioVegng7Hfx0XmIUavZgz5UL-HIgPCBECc-Ws=s128",
		"SHIB": "https://mixin-images.zeromesh.net/fgSEd6CY07BiZP76--7JA9P-rKIWRoXD8Eis8RUL6mP85_QPsbMoyJtWJ6MjE9jWFEjabNF0AKb8i2QOfdbCS6BJMntySps-8GfvJQ=s128",
	}

	SourceTitle = map[string]string{
		core.TransferSourceTradeConfirmed:           TitleTradeConfirmed,
		core.TransferSourceOrderCancelled:           TitleOrderCancelled,
		core.TransferSourceOrderFilled:              TitleOrderFilled,
		core.TransferSourceOrderInvalid:             TitleOrderInvalid,
		core.TransferSourcePositionExercised:        TitlePositionExercised,
		core.TransferSourcePositionExercisedRefund:  TitlePositionExercisedRefund,
		core.TransferSourcePositionExercisedInvalid: TitlePositionExercisedInvalid,
		core.TransferSourcePositionClosedRefund:     TitlePositionClosedRefund,
	}
)

func BuildTransferCardMsg(ctx context.Context, userId, amount, assetId, traceId string) (m core.Message, err error) {
	messenger := config.Cfg.DApp.AppID
	conversationId := bot.UniqueConversationId(messenger, userId)
	action := fmt.Sprintf("mixin://snapshots?trace=%s", traceId)
	currency := core.GetCurrencyByAsset(assetId, false)

	data, _ := json.Marshal(map[string]string{
		"title":       amount,
		"description": currency,
		"action":      action,
		"icon_url":    iconMaps[currency],
		"app_id":      messenger,
	})

	msgReq := mixin.MessageRequest{
		ConversationID: conversationId,
		RecipientID:    userId,
		MessageID:      traceId,
		Category:       MessageCategoryAppCard,
		Data:           base64.StdEncoding.EncodeToString(data),
	}
	return core.ToMessage(&msgReq)
}

func BuildBizInfoMsg(ctx context.Context, source, userId string, transfer *core.Transfer) (m core.Message, err error) {
	messenger := config.Cfg.DApp.AppID
	msgContent := BuildMsgContent(ctx, source, userId, transfer)
	messageId := core.GetSettlementId(transfer.TransferId, transfer.Source)
	msgReq := mixin.MessageRequest{
		ConversationID: bot.UniqueConversationId(messenger, userId),
		RecipientID:    userId,
		MessageID:      messageId,
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msgContent)),
	}
	return core.ToMessage(&msgReq)
}

func BuildExerciseNotifyMsg(ctx context.Context, userId, positionId, content, modifier string) (m core.Message, err error) {
	messenger := config.Cfg.DApp.AppID
	messageId := core.GetSettlementId(positionId, modifier)
	msgReq := mixin.MessageRequest{
		ConversationID: bot.UniqueConversationId(messenger, userId),
		RecipientID:    userId,
		MessageID:      messageId,
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(content)),
	}
	return core.ToMessage(&msgReq)
}

func GetBizInfoMsgByTrade(ctx context.Context, source string, trade *core.Trade, transfer *core.Transfer) (msg core.Message, err error) {

	messenger := config.Cfg.DApp.AppID
	userId := transfer.UserId
	tradeAmount := number.FromString(trade.Price).Mul(number.FromString(trade.Amount)).String()
	msgContent := BuildMsg(source, trade.InstrumentName, trade.Side, trade.Amount, trade.Amount, "", tradeAmount, trade.BaseAssetID, trade.QuoteAssetID)
	messageId := core.GetSettlementId(transfer.TransferId, transfer.Source)

	msgReq := mixin.MessageRequest{
		ConversationID: bot.UniqueConversationId(messenger, userId),
		RecipientID:    userId,
		MessageID:      messageId,
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msgContent)),
	}
	return core.ToMessage(&msgReq)
}

func GetExpiryNotifyContent(instrumentName, side string, size float64) (res string, err error) {
	i, err := core.ParseInstrument(instrumentName)
	if err != nil {
		return "", err
	}
	dateCn := i.ExpirationDate.Format("2006年01月02日")
	sizeStr := decimal.NewFromFloat(size).String()
	return fmt.Sprintf(TpltExpiryNotify, dateCn, i.StrikePrice, i.BaseCurrency, sizeStr, i.BaseCurrency), nil
}

func BuildMsgContent(ctx context.Context, source, userId string, transfer *core.Transfer) string {
	var msgContent string
	switch source {
	case core.TransferSourceOrderCancelled, core.TransferSourceOrderInvalid:
		var order core.Order
		if err := config.Db.Model(core.Order{}).Where("order_id=?", transfer.Detail).Find(&order).Error; err != nil {
			break
		}
		totalSize, filledSize, remainingSize := GetOrderSize(order)
		msgContent = BuildMsg(source, order.InstrumentName, order.Side, totalSize, filledSize, remainingSize, transfer.Amount, order.BaseAssetID, order.QuoteAssetID)
		break
	case core.TransferSourceTradeConfirmed:
		var trade core.Trade
		if err := config.Db.Model(core.Trade{}).Where("trade_id=? AND user_id = ?", transfer.Detail, transfer.UserId).Find(&trade).Error; err != nil {
			break
		}
		//var orderID string
		//var order p.Order
		//if len(trade.Side) > 0 {
		//	if trade.Side == "ASK" {
		//		orderID = trade.AskOrderID
		//	}
		//	if trade.Side == "BID" {
		//		orderID = trade.BidOrderID
		//	}
		//}
		//if err := p.OrderMgr(core.Db).Where("order_id=?", orderID).Find(&order).Error; err != nil {
		//	break
		//}
		//totalSize, filledSize := GetOrderFilledSize(order)
		msgContent = BuildMsg(source, trade.InstrumentName, trade.Side, trade.Amount, trade.Amount, "", transfer.Amount, trade.BaseAssetID, trade.QuoteAssetID)
		break
	case core.TransferSourcePositionExercised, core.TransferSourcePositionExercisedInvalid, core.TransferSourcePositionExercisedRefund:
		var position core.Position
		if err := config.Db.Model(core.Position{}).Where("position_id = ?", transfer.Detail).Find(&position).Error; err != nil {
			break
		}
		if len(position.PositionID) > 0 {
			totalSize, filledSize := number.FromFloat(position.Size).String(), number.FromFloat(position.ExercisedSize).String()
			msgContent = BuildMsg(source, position.InstrumentName, position.Side, totalSize, filledSize, "", transfer.Amount, transfer.AssetId, transfer.AssetId)
		}
		break
	case core.TransferSourcePositionClosedRefund:
		var trade core.Trade
		split := strings.Split(transfer.Detail, ":")
		if err := config.Db.Model(core.Trade{}).Where("trade_id = ? AND user_id = ? AND status = 10", split[1], transfer.UserId).Find(&trade).Error; err != nil {
			break
		}
		if len(trade.TradeID) > 0 {
			msgContent = BuildMsg(source, trade.InstrumentName, trade.Side, trade.Amount, trade.Amount, "", transfer.Amount, trade.BaseAssetID, trade.QuoteAssetID)
		}
		break
	}
	return msgContent
}

func GetOrderSize(order core.Order) (totalSize, filledSize, remainingSize string) {
	if order.Side == "ASK" {
		filledSize = order.FilledAmount
		remainingSize = order.RemainingAmount
		totalSize = number.FromString(order.FilledAmount).Add(number.FromString(order.RemainingAmount)).String()
	}
	if order.Side == "BID" {
		filledSize = order.FilledAmount
		remainingSize = number.FromString(order.RemainingFunds).Div(number.FromString(order.Price)).String()
		totalSize = number.FromString(order.FilledFunds).
			Add(number.FromString(order.RemainingFunds)).
			Div(number.FromString(order.Price)).String()
	}
	return
}

func BuildMsg(source, instrument, side, totalSize, filledSize, remainingSize, amount, baseAssetId, quoteAssetId string) string {
	if instrument == "" {
		return ""
	}
	var builder = &TransferMessageBuilder{}
	return builder.New().SetSource(source).
		SetTitle(SourceTitle[source]).
		SetInstrument(side, instrument).
		SetSide(side, builder.msg.OptionType).
		SetSize(side, totalSize, filledSize, remainingSize, baseAssetId, quoteAssetId).
		SetAmount(side, amount, baseAssetId, quoteAssetId).
		BuildMessage()
}

func UniqueConversationId(userId, recipientId string) string {
	minId, maxId := userId, recipientId
	if strings.Compare(userId, recipientId) > 0 {
		maxId, minId = userId, recipientId
	}
	h := md5.New()
	io.WriteString(h, minId)
	io.WriteString(h, maxId)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}

package message

import (
	"context"
	"encoding/base64"
	"github.com/fox-one/mixin-sdk-go"
	"math/rand"
	"option-dance/cmd/config"
	"option-dance/core"
	"strconv"
)

type messageBuilder struct {
}

func NewMessageBuilder() core.MessageBuilder {
	return &messageBuilder{}
}

func (m *messageBuilder) OrderCreateGroupNotify(ctx context.Context, order *core.Order) (*core.Message, error) {

	tplt, err := GetTplt(TpltCreateOrderGroupNtf, order)
	if err != nil {
		return nil, err
	}
	modifier := rand.Intn(1)
	createOrderNtfMessage, err := core.ToMessage(&mixin.MessageRequest{
		ConversationID: config.Cfg.DApp.GroupConversationId,
		RecipientID:    "",
		MessageID:      core.GetSettlementId(order.OrderID, strconv.Itoa(modifier)),
		Category:       MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(tplt)),
	})
	if err != nil {
		return nil, err
	}
	return &createOrderNtfMessage, nil
}

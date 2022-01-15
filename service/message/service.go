package message

import (
	"context"
	"encoding/json"
	"github.com/fox-one/mixin-sdk-go"
	"go.uber.org/zap"
	"option-dance/core"
)

func New(client *mixin.Client, messages core.MessageStore, positionStore core.PositionStore) core.MessageService {
	return &messageService{c: client, messageStore: messages, positionStore: positionStore}
}

type messageService struct {
	c             *mixin.Client
	messageStore  core.MessageStore
	positionStore core.PositionStore
}

func (s *messageService) Send(ctx context.Context, message *core.Message) error {
	raw := json.RawMessage(message.Data)
	err := s.c.SendRawMessage(ctx, raw)
	zap.L().Debug("send message:", zap.Any("raw", raw))
	if mixin.IsErrorCodes(err, 10002) {
		return nil
	}
	return err
}

func (s *messageService) SendBatch(ctx context.Context, messages []*core.Message) error {
	raws := make([]json.RawMessage, 0, len(messages))
	for _, msg := range messages {
		raws = append(raws, json.RawMessage(msg.Data))
	}
	err := s.c.SendRawMessages(ctx, raws)
	zap.L().Debug("batch send messages:", zap.Any("raws", raws))
	if mixin.IsErrorCodes(err, 10002) {
		return nil
	}
	return err
}

func (s *messageService) Meet(ctx context.Context, userID string) error {
	_, err := s.c.CreateContactConversation(ctx, userID)
	return err
}

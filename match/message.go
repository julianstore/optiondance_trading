package match

import (
	"context"
	"go.uber.org/zap"
	"option-dance/core"
	"time"
)

func (ex *Exchange) LoopingSendMessage(ctx context.Context) error {
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := ex.SendMessage(ctx); err == nil {
				dur = 200 * time.Millisecond
			} else {
				zap.L().Info("LoopingSendMessage error: ", zap.Error(err))
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (ex *Exchange) SendMessage(ctx context.Context) error {
	const Limit = 300
	const Batch = 70

	messages, err := ex.messageStore.List(ctx, Limit)
	if err != nil {
		zap.L().Error("messageStore.List", zap.Error(err))
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	filter := make(map[string]bool)
	var idx int

	for _, msg := range messages {
		if filter[msg.UserID] {
			continue
		}

		messages[idx] = msg
		filter[msg.UserID] = true
		idx++

		if idx >= Batch {
			break
		}
	}

	messages = messages[:idx]

	var p2pMessages = make([]*core.Message, 0)
	var groupMessages = make([]*core.Message, 0)

	for _, e := range messages {
		if e.UserID == "" {
			groupMessages = append(groupMessages, e)
		} else {
			p2pMessages = append(p2pMessages, e)
		}
	}

	if err := ex.messageService.SendBatch(ctx, p2pMessages); err != nil {
		zap.L().Error("messageService.SendBatch", zap.Error(err))
		return err
	}

	for _, e := range groupMessages {
		if err := ex.messageService.Send(ctx, e); err != nil {
			zap.L().Error("messageService.SendBatch", zap.Error(err))
			return err
		}
	}

	if err := ex.messageStore.Delete(ctx, messages); err != nil {
		zap.L().Error("messageStore.Delete", zap.Error(err))
		return err
	}

	for _, message := range messages {
		zap.L().Debug("transfer msg sended", zap.Any("msg", message))
	}
	return nil
}

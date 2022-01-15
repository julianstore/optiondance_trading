package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fox-one/mixin-sdk-go"
	"time"
)

type (
	Message struct {
		ID        int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		MessageID string    `gorm:"column:message_id;type:varchar(45)" json:"message_id"`
		UserID    string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		Data      string    `gorm:"column:data;type:text" json:"data"`
		CreatedAt time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
	}

	MessageRequest struct {
		ConversationID   string `json:"conversation_id"`
		RecipientID      string `json:"recipient_id,omitempty"`
		MessageID        string `json:"message_id"`
		Category         string `json:"category"`
		Data             string `json:"data"`
		RepresentativeID string `json:"representative_id,omitempty"`
		QuoteMessageID   string `json:"quote_message_id,omitempty"`
	}

	MessageStore interface {
		Create(ctx context.Context, message *Message) error
		CreateBatch(ctx context.Context, messages []*Message) error
		List(ctx context.Context, limit int) ([]*Message, error)
		Delete(ctx context.Context, messages []*Message) error
	}

	MessageService interface {
		Send(ctx context.Context, message *Message) error
		SendBatch(ctx context.Context, messages []*Message) error
		Meet(ctx context.Context, userID string) error
	}

	MessageBuilder interface {
		OrderCreateGroupNotify(ctx context.Context, order *Order) (*Message, error)
	}
)

func ToMessage(req *mixin.MessageRequest) (message Message, err error) {
	data, _ := json.Marshal(&req)
	if req.Data == "" {
		return Message{}, fmt.Errorf("empty message content")
	}
	return Message{
		MessageID: req.MessageID,
		UserID:    req.RecipientID,
		Data:      string(data),
		CreatedAt: time.Now(),
	}, nil
}

package core

import (
	"time"
)

type (
	MixinMessage struct {
		ID               int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		ConversationID   string    `gorm:"column:conversation_id;type:varchar(45)" json:"conversation_id"`
		UserID           string    `gorm:"column:user_id;type:varchar(45)" json:"user_id"`
		MessageID        string    `gorm:"column:message_id;type:varchar(45)" json:"message_id"`
		Category         string    `gorm:"column:category;type:varchar(45)" json:"category"`
		Data             string    `gorm:"column:data;type:text" json:"data"`
		DecodeData       string    `gorm:"column:decode_data;type:text" json:"decode_data"`
		RepresentativeID string    `gorm:"column:representative_id;type:varchar(45)" json:"representative_id"`
		QuoteMessageID   string    `gorm:"column:quote_message_id;type:varchar(45)" json:"quote_message_id"`
		Status           string    `gorm:"column:status;type:varchar(45)" json:"status"`
		Source           string    `gorm:"column:source;type:varchar(45)" json:"source"`
		CreatedAt        time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
		UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	}

	MixinUser struct {
		ID             int64     `gorm:"primaryKey;autoIncrement;not null;" json:"-"`
		UserID         string    `gorm:"column:user_id;type:varchar(64)" json:"user_id"`
		Phone          string    `gorm:"column:phone;type:varchar(45)" json:"phone"`
		FullName       string    `gorm:"column:full_name;type:varchar(45)" json:"full_name"`
		AvatarURL      string    `gorm:"column:avatar_url;type:varchar(100)" json:"avatar_url"`
		IdentityNumber string    `gorm:"column:identity_number;type:varchar(45)" json:"identity_number"`
		SessionID      string    `gorm:"column:session_id;type:varchar(64)" json:"session_id"`
		PinToken       string    `gorm:"column:pin_token;type:varchar(64)" json:"pin_token"`
		DeviceStatus   string    `gorm:"column:device_status;type:varchar(45)" json:"device_status"`
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
		UpdateTime     time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
		Status         int8      `gorm:"column:status;type:tinyint" json:"status"`
	}
)

func (m *MixinMessage) TableName() string {
	return "mixin_message"
}

func (m *MixinUser) TableName() string {
	return "mixin_user"
}

type (
	MultisigDTO struct {
		AssetId          string           `json:"asset_id"`
		Amount           string           `json:"amount"`
		TraceId          string           `json:"trace_id"`
		Memo             string           `json:"memo"`
		OpponentMultisig OpponentMultisig `json:"opponent_multisig"`
	}

	OpponentMultisig struct {
		Receivers []string `json:"receivers"`
		Threshold int64    `json:"threshold"`
	}

	MultisigResponseDTO struct {
		Data MultisigResponse `json:"data"`
	}

	MultisigResponse struct {
		Type      string    `json:"type"`
		TraceId   string    `json:"trace_id"`
		AssetId   string    `json:"asset_id"`
		Amount    string    `json:"amount"`
		Threshold int64     `json:"threshold"`
		Receivers []string  `json:"receivers"`
		Memo      string    `json:"memo"`
		CreatedAt time.Time `json:"created_at"`
		Status    string    `json:"status"`
		CodeId    string    `json:"code_id"`
	}

	MixinUserEntity struct {
		Type                     string    `json:"type"`
		UserID                   string    `json:"user_id"`
		IdentityNumber           string    `json:"identity_number"`
		Phone                    string    `json:"phone"`
		FullName                 string    `json:"full_name"`
		Biography                string    `json:"biography"`
		AvatarURL                string    `json:"avatar_url"`
		Relationship             string    `json:"relationship"`
		MuteUntil                time.Time `json:"mute_until"`
		CreatedAt                time.Time `json:"created_at"`
		IsVerified               bool      `json:"is_verified"`
		IsScam                   bool      `json:"is_scam"`
		SessionID                string    `json:"session_id"`
		PinToken                 string    `json:"pin_token"`
		PinTokenBase64           string    `json:"pin_token_base64"`
		CodeID                   string    `json:"code_id"`
		CodeURL                  string    `json:"code_url"`
		DeviceStatus             string    `json:"device_status"`
		HasPin                   bool      `json:"has_pin"`
		HasEmergencyContact      bool      `json:"has_emergency_contact"`
		ReceiveMessageSource     string    `json:"receive_message_source"`
		AcceptConversationSource string    `json:"accept_conversation_source"`
		AcceptSearchSource       string    `json:"accept_search_source"`
		FiatCurrency             string    `json:"fiat_currency"`
		//TransferNotificationThreshold int       `json:"transfer_notification_threshold"`
		//TransferConfirmationThreshold int       `json:"transfer_confirmation_threshold"`
	}

	MixinUserResponse struct {
		Data MixinUserEntity `json:"data"`
	}

	SettingsDTO struct {
		AppMode      int8 `json:"app_mode"`
		DeliveryType int8 `json:"delivery_type"`
	}
)

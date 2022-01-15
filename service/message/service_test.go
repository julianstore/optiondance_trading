package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/MixinNetwork/bot-api-go-client"
	msdk "github.com/fox-one/mixin-sdk-go"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/pkg/mixin"
	"testing"
	"time"
)

var ctx = context.Background()
var cfg config.Config
var client *msdk.Client

const odDevMarketMakerGroupConversationID = "252be0c0-0b7d-4bd7-a5e8-0a12e7caf3b0"

func TestMain(m *testing.M) {
	config.InitTestConfig("prod")
	cfg = config.Cfg
	var err error
	client, err = mixin.Client()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestMessageService_Send(t *testing.T) {
	var msgList = make([]*core.Message, 0)
	mid := "222bd269-0a1a-372f-8ac2-36cedc43ba0d"
	t.Log(mid)
	messenger := config.Cfg.DApp.AppID
	userId := "019308a5-e1a9-427c-af2a-e05093beedaa"
	messageView := msdk.MessageRequest{
		ConversationID: bot.UniqueConversationId(messenger, userId),
		RecipientID:    userId,
		MessageID:      mid,
		Category:       msdk.MessageCategoryAppCard,
		Data:           "eyJhY3Rpb24iOiJtaXhpbjovL3NuYXBzaG90cz90cmFjZT0yMjJiZDI2OS0wYTFhLTM3MmYtOGFjMi0zNmNlZGM0M2JhMGIiLCJhcHBfaWQiOiIzNTEzMGYyNS00MzVhLTRiYTEtOGEyOC0xZDgxNjE4MWI3ZTkiLCJkZXNjcmlwdGlvbiI6IkNOQiIsImljb25fdXJsIjoiaHR0cHM6Ly9taXhpbi1pbWFnZXMuemVyb21lc2gubmV0LzBzUVk2M2RETWtXVFVSa0pWam93V1k2TGU0SUNqQUZ1dTNBTlZ5WkE0dUkzVWRrYnVPVDVmakpVVDgyQXJOWW1adlZjeERYeU5qeG9PdjBUQVliUVROS1M9czEyOCIsInRpdGxlIjoiNDAwMCJ9",
	}
	marshal, err := json.Marshal(&messageView)
	if err != nil {
		t.Error(err)
	}
	msgList = append(msgList, &core.Message{
		MessageID: mid,
		UserID:    userId,
		Data:      string(marshal),
		CreatedAt: time.Now(),
	})
	//client, _ := mixin.Client()
	//messagez := New(client, provider.ProvideMessageStore())
	//err = messagez.SendBatch(context.Background(), msgList)
	if err != nil {
		t.Error(err)
	}
}

//  OptionDance Market Maker dev Group conversationID: 252be0c0-0b7d-4bd7-a5e8-0a12e7caf3b0
//  OptionDance Market Maker beta Group conversationID: 13b651a3-f686-4904-998a-e6a91d8b33f3
//  OptionDance Market Maker uat Group conversationID: c08f26fe-4df3-4dd3-b208-fe5be809efa7
//  OptionDance Market Maker prod Group conversationID: eb8a7f9f-c4be-4437-8f65-cff5534fe808
func TestCreateGroup(t *testing.T) {
	client, err := mixin.Client()
	conversationID := bot.UuidNewV4().String()
	participants := make([]*msdk.Participant, 0)
	createdAt := time.Now()
	participants = append(participants,
		&msdk.Participant{Type: "participant", UserID: config.Cfg.DApp.AppID, Role: "OWNER", CreatedAt: createdAt},
		&msdk.Participant{Type: "participant", UserID: "019308a5-e1a9-427c-af2a-e05093beedaa", Role: "", CreatedAt: createdAt},
		//&msdk.Participant{Type: "participant", UserID: "5d0b24d1-f2f5-4264-8abc-cdb51b18ee50", Role: "", CreatedAt: createdAt},
	)
	conversation, err := client.CreateGroupConversation(ctx, conversationID, "OptionDance Market Maker Group", participants)
	if err != nil {
		t.Error(err)
	}
	t.Log(conversationID, conversation)
}

func TestSendGroupMessage(t *testing.T) {
	client, err := mixin.Client()
	if err != nil {
		t.Error(err)
	}
	//err = client.SendMessage(ctx, &msdk.MessageRequest{
	//	ConversationID:   odDevMarketMakerGroupConversationID,
	//	RecipientID:      "",
	//	MessageID:        bot.UuidNewV4().String(),
	//	Category:         msdk.MessageCategoryPlainText,
	//	Data:             base64.StdEncoding.EncodeToString([]byte("你好")),
	//	RepresentativeID: "",
	//	QuoteMessageID:   "",
	//})

	err = client.SendMessages(ctx, []*msdk.MessageRequest{
		{

			ConversationID:   odDevMarketMakerGroupConversationID,
			RecipientID:      "",
			MessageID:        bot.UuidNewV4().String(),
			Category:         msdk.MessageCategoryPlainText,
			Data:             base64.StdEncoding.EncodeToString([]byte("你好")),
			RepresentativeID: "",
			QuoteMessageID:   "",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAddMember(t *testing.T) {
	_, err := client.AddParticipants(ctx, cfg.DApp.GroupConversationId, "443bf7c0-308f-461e-ac0f-88ef19b9162b")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveMember(t *testing.T) {
	_, err := client.RemoveParticipants(ctx, cfg.DApp.GroupConversationId, "443bf7c0-308f-461e-ac0f-88ef19b9162b")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateGroupName(t *testing.T) {
	_, err := client.UpdateConversation(ctx, cfg.DApp.GroupConversationId, msdk.ConversationUpdate{
		Name:         "OptionDance.DEV market maker group",
		Announcement: "test Announcement",
	})
	if err != nil {
		t.Error(err)
	}
}

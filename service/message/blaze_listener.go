package message

import (
	"context"
	"encoding/json"
	"option-dance/cmd/config"
	"option-dance/pkg/mixin"

	"go.uber.org/zap"
)

var (
	Blaze        *BlazeClient
	uid, sid, pk string
)

func MixinMessageListener() {
	defer MixinMessageListener()
	zap.L().Info("Start MixinFollowMsgListenser")
	m := config.Cfg.DApp
	uid, sid, pk = m.AppID, m.SessionID, m.PrivateKey
	Blaze = NewBlazeClient(uid, sid, pk)
	var l MixinFollowMsgListenser
	err := Blaze.Loop(context.Background(), l)
	if err != nil {
		zap.L().Error("Msg listener error:" + err.Error())
	}
	zap.L().Info("Stop MixinFollowMsgListenser")
}

type MixinFollowMsgListenser struct {
}

func (l MixinFollowMsgListenser) OnMessage(ctx context.Context, msg MessageView, userId string) error {
	MessageReaded(msg.MessageId) //Message has been read
	return nil
}

type MixinMsgListenser struct{}

func (l MixinMsgListenser) OnMessage(ctx context.Context, msg MessageView, userId string) error {
	return nil
}

//Robot message has been read
func MessageReaded(messageIds ...string) (success bool) {
	type MessageStatus struct {
		MessageId string `json:"message_id"`
		Status    string `json:"status"`
	}
	var statusList []MessageStatus
	for _, e := range messageIds {
		statusList = append(statusList, MessageStatus{
			MessageId: e,
			Status:    "READ",
		})
	}
	postData, err := json.Marshal(statusList)
	println(string(postData))
	if err != nil {
		zap.L().Error(err.Error())
	}
	_, err = mixin.MixinRequest("/acknowledgements", "POST", postData)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return
}

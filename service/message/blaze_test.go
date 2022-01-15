package message

import (
	"context"
	"fmt"
	"option-dance/cmd/config"
	"testing"
)

func TestSendAppCardMsg(t *testing.T) {
	currency := "CNB"
	clientId := config.Cfg.DApp.AppID
	m := config.Cfg.DApp
	uid, sid, pk := m.AppID, m.SessionID, m.PrivateKey
	blaze := NewBlazeClient(uid, sid, pk)
	var l MixinFollowMsgListenser
	go blaze.Loop(context.Background(), l)
	userId := "019308a5-e1a9-427c-af2a-e05093beedaa"
	conversationId := UniqueConversationId(clientId, userId)
	action := fmt.Sprintf("mixin://snapshots?trace=%s", "85c6eef3-bc72-3575-9502-b2b2d50c581c")
	iconUrl := "https://mixin-images.zeromesh.net/0sQY63dDMkWTURkJVjowWY6Le4ICjAFuu3ANVyZA4uI3UdkbuOT5fjJUT82ArNYmZvVcxDXyNjxoOv0TAYbQTNKS=s128"
	blaze.SendAppCard(context.Background(), conversationId, userId, "98812", currency, action, m.AppID, iconUrl)
}

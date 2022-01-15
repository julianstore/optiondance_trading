package mixin

import (
	"bytes"
	"github.com/MixinNetwork/bot-api-go-client"
	mixinSdk "github.com/fox-one/mixin-sdk-go"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"option-dance/cmd/config"
	"strings"
)

const (
	uri = "https://api.mixin.one"
)

func InitBotParams() (uid string, sid string, pk string, pin string, pinToken string, assetId string) {
	m := config.Cfg.DApp
	uid, sid, pk, pin, pinToken = m.AppID, m.SessionID, m.PrivateKey, m.Pin, m.PinToken
	return
}

func MixinRequest(path string, method string, body []byte) (res []byte, err error) {
	uid, sid, pk, _, _, _ := InitBotParams()
	var bodyStr string = ""
	if body != nil && len(body) > 0 {
		bodyStr = string(body)
	}
	token, err := bot.SignAuthenticationToken(uid, sid, pk, method, path, bodyStr)
	if err != nil {
		println(err)
		return
	}
	req, err := http.NewRequest(method, uri+path, bytes.NewBuffer(body))
	if err != nil {
		println(err)
		return
	}

	httpClient := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-Id", uuid.NewV4().String())
	resp, err := httpClient.Do(req)
	if err != nil {
		println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		println(err)
		return
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

func PayResultByTraceIdUseRobot(id string) (success bool) {
	res, err := MixinRequest("/transfers/trace/"+id, "GET", nil)
	if strings.Contains(string(res), "error") || err != nil {
		success = false
	} else {
		success = true
	}
	return
}

func Client() (*mixinSdk.Client, error) {
	clientId, sid, pk, _, token, _ := InitBotParams()
	keystore := &mixinSdk.Keystore{
		ClientID:   clientId,
		SessionID:  sid,
		PrivateKey: pk,
		PinToken:   token,
	}
	return mixinSdk.NewFromKeystore(keystore)
}

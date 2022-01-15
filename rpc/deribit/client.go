package deribit

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"option-dance/cmd/config"
	"strconv"
	"strings"
)

const (
	DeribitURL = "https://www.deribit.com"
)

type Client struct {
	URL           string
	Client        *http.Client
	VersionPrefix string
	ClientId      string
	Secret        string
}

func New(url, versionPrefix, clientId, secret string) *Client {
	return &Client{
		URL:           url,
		Client:        &http.Client{},
		VersionPrefix: versionPrefix,
		ClientId:      clientId,
		Secret:        secret,
	}
}

func C() *Client {
	deribit := config.Cfg.Deribit
	return New(DeribitURL, "/api/v2", deribit.ClientId, deribit.Secret)
}

func (c *Client) Get(path string, query map[string]interface{}) ([]byte, error) {
	url := c.URL + c.VersionPrefix + path + "?" + toQueryString(query)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "private") {
		s := c.ClientId + ":" + c.Secret
		authorization := "Basic " + base64.URLEncoding.EncodeToString([]byte(s))
		request.Header.Add("Authorization", authorization)
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	//util.JsonPrintS(string(body))
	return body, nil
}

func toQueryString(queryMap map[string]interface{}) (qs string) {
	if queryMap == nil {
		return ""
	}
	for k, v := range queryMap {
		var value string
		i, ok := v.(int)
		if ok {
			value = strconv.Itoa(i)
		}
		switch v.(type) {
		case int:
			value = strconv.Itoa(i)
			break
		case string:
			value = fmt.Sprintf("%v", v)
		}
		if len(value) == 0 {
			continue
		}
		qs += k + "=" + value + "&"
	}
	qs = qs[:len(qs)-1]
	return
}

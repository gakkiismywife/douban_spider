package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"spider_douban/cache"
	"spider_douban/config"
	"strings"
	"time"
)

type MessageRequest struct {
	ToUser      string            `json:"touser"`
	MessageType string            `json:"msgtype"`
	AgentId     string            `json:"agentid"`
	Text        map[string]string `json:"text"`
}

type TokenResponse struct {
	Code    int    `json:"errcode"`
	Msg     string `json:"errormsg"`
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

func NewMessageRequest(content string) *MessageRequest {
	m := make(map[string]string)
	m["content"] = content
	return &MessageRequest{
		ToUser:      "@all",
		MessageType: "text",
		AgentId:     config.Wechat.AgentId,
		Text:        m,
	}
}

func SendMessage(token, content string) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)

	client := &http.Client{}

	reqParam := NewMessageRequest(content)

	reqByte, err := json.Marshal(reqParam)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}

	body := strings.NewReader(string(reqByte))

	req, _ := http.NewRequest("POST", url, body)

	do, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do err", err)
		return
	}

	defer do.Body.Close()

	all, _ := ioutil.ReadAll(do.Body)
	fmt.Println(string(all))
}

// GetAccessToken 获取token
func GetAccessToken(id, secret string) string {
	rdb := cache.GetRedisClient()
	defer rdb.Close()
	result, _ := rdb.Get(context.Background(), "access_token").Result()
	if result != "" {
		return result
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", id, secret)
	response, err := http.Get(url)
	if err != nil {
		log.Println("[GetAccessToken] failed,err:", err.Error())
		return ""
	}
	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	res := new(TokenResponse)
	err = json.Unmarshal(all, res)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return ""
	}

	rdb.Set(context.Background(), "access_token", res.Token, time.Second*1800)

	return res.Token
}

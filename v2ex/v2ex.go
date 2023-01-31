package v2ex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-module/carbon/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"spider_douban/cache"
	"spider_douban/config"
	"spider_douban/db"
	"spider_douban/wechat"
)

type Item struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Created int64  `json:"created"`
}

type Response struct {
	Result []Item
}

func (i *Item) GenerateMessage(title, url, time string) string {
	return fmt.Sprintf("监测到新的二手交易贴\n标题：%s\n链接：%s\n发布时间：%s", title, url, time)
}

func Run() {

	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	request, _ := http.NewRequest("GET", config.V2ex.Url, nil)

	request.Header.Add("Authorization", config.V2ex.Token)

	response, err := client.Do(request)
	if err != nil {
		log.Printf("[v2ex]request failed,err : %v", err.Error())
		return
	}

	defer response.Body.Close()

	resp := new(Response)

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[v2ex]read failed,err : %v", err.Error())
		return
	}

	err = json.Unmarshal(all, resp)
	if err != nil {
		log.Printf("[v2ex]json decode failed,err : %v", err.Error())
		return
	}

	rc := cache.GetRedisClient()
	defer rc.Close()

	last, _ := rc.Get(context.Background(), "v2ex_last_request").Int64()
	wechatToken := wechat.GetAccessToken(config.Wechat.Key, config.Wechat.Secret)

	now := carbon.Now().Timestamp()
	for _, i := range resp.Result {
		if last != 0 && last > now {
			log.Println(fmt.Sprintf("[v2ex]old topic ,name:%s,url:%s,time:%s", i.Title, i.Url, carbon.CreateFromTimestamp(i.Created).ToDateTimeString()))
			continue
		}
		db.CreateMessage(i.Title, i.Url)
		t := carbon.CreateFromTimestamp(i.Created)
		wechat.SendMessage(wechatToken, i.GenerateMessage(i.Title, i.Url, t.ToDateTimeString()))
	}

	//更新上次请求时间
	rc.Set(context.Background(), "v2ex_last_request", now, redis.KeepTTL)
}

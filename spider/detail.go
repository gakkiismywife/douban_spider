package spider

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"spider_douban/cache"
	"spider_douban/config"
	"spider_douban/wechat"
	"strings"
	"time"
)

type pageTask struct {
	*Task
	Title string
}

func NewPageTask(flag, url, title string) *pageTask {
	t := NewTask(flag, url)

	p := &pageTask{
		Task:  nil,
		Title: title,
	}
	t.SetHtmlCallback(".create-time.color-green", p.htmlHandle)
	t.SetRequestCallback(p.requestHandle)
	t.SetResponseCallback(p.responseHandle)
	p.Task = t
	return p
}

func (p *pageTask) htmlHandle(e *colly.HTMLElement) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Text, time.Local)
	publishTime := t.Unix()
	log.Println(fmt.Sprintf("[%s]%s创建时间为%s", p.Flag, p.Title, e.Text))

	rdb := cache.GetRedisClient()
	defer rdb.Close()
	go rdb.HSet(context.Background(), config.Task.Home, p.Url, p.Title).Result() //放入缓存

	if time.Now().Unix()-publishTime < (config.INTERVAL * 2) {
		go send(p.Title, p.Url, e.Text)
	}
}

func (p *pageTask) requestHandle(request *colly.Request) {
	request.Headers.Set("Referer", config.Task.Home)
	log.Println(fmt.Sprintf("[%s]Visiting:%s", p.Flag, p.Title))
}

func (p *pageTask) responseHandle(response *colly.Response) {
	body := string(response.Body)
	//判断响应是否正常
	if !strings.Contains(body, "create-time") {
		log.Println(fmt.Sprintf("[%s]response body err:%s", p.Flag, body))
		p.State = false
	} else {
		p.State = true
	}
}

func send(title, url, publishTime string) {
	message := fmt.Sprintf("监测到新的帖子\n标题：%s\n链接：%s\n发布时间：%s", title, url, publishTime)
	token := wechat.GetAccessToken(config.Wechat.Key, config.Wechat.Secret)
	if token == "" {
		log.Println("[process]获取token失败")
		return
	}
	wechat.SendMessage(token, message)
}

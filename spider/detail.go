package spider

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"spider_douban/cache"
	"spider_douban/config"
	"spider_douban/db"
	"spider_douban/wechat"
	"strings"
	"time"
)

type pageTask struct {
	*Task
	Title string
}

func (p *pageTask) GenerateMessage(title, url, time string) string {
	return fmt.Sprintf("监测到新的帖子\n标题：%s\n链接：%s\n发布时间：%s", title, url, time)
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
	t.SetErrorCallback(p.errorHandle)
	p.Task = t
	return p
}

func (p *pageTask) htmlHandle(e *colly.HTMLElement) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Text, time.Local)
	publishTime := t.Unix()
	log.Println(fmt.Sprintf("[%s]%s 创建时间为%s", p.Flag, p.Title, e.Text))

	if time.Now().Unix()-publishTime < int64(config.Task.Seconds) {
		go p.Send(e.Text)
	}
}

func (p *pageTask) requestHandle(request *colly.Request) {
	request.Headers.Set("Referer", config.Task.Home)
	log.Println(fmt.Sprintf("[%s] Visiting:%s", p.Flag, p.Title))
}

func (p *pageTask) responseHandle(response *colly.Response) {
	body := string(response.Body)
	//判断响应是否正常
	if !strings.Contains(body, "create-time") {
		log.Println(fmt.Sprintf("[%s]%sResponse body err ", p.Flag, p.Title))
		p.State = false
		rbd := cache.GetRedisClient()
		defer rbd.Close()
		rbd.HDel(context.Background(), config.Task.Home, p.Url)
	} else {
		p.State = true
	}
}

func (p *pageTask) errorHandle(response *colly.Response, err error) {
	rbd := cache.GetRedisClient()
	defer rbd.Close()
	rbd.HDel(context.Background(), config.Task.Home, p.Url)
}

func (p *pageTask) Send(publishTime string) {
	if db.HasSend(p.Title, p.Url) {
		log.Println(fmt.Sprintf("[%s] %s已经发送过消息", p.Flag, p.Title))
		return
	}
	message := p.GenerateMessage(p.Title, p.Url, publishTime)
	token := wechat.GetAccessToken(config.Wechat.Key, config.Wechat.Secret)
	if token == "" {
		log.Println(fmt.Sprintf("[%s]获取token失败", p.Flag))
		return
	}
	db.CreateMessage(p.Title, p.Url)
	wechat.SendMessage(token, message)
}

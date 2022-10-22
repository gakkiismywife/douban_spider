package process

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"spider_douban/cache"
	"spider_douban/config"
	"spider_douban/ip"
	"spider_douban/wechat"
	"time"
)

// VisitDetail 请求帖子详情
func VisitDetail(detailUrl, title string) {
	rdb := cache.GetRedisClient()
	defer rdb.Close()

	//判断是否请求过
	isVisited := rdb.HExists(context.Background(), config.Task.Home, detailUrl).Val()
	if isVisited {
		fmt.Println(fmt.Sprintf("[process]%s has visited", title))
		return
	}
	c := initCollector()

	c.OnHTML(".create-time.color-green", func(e *colly.HTMLElement) {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Text, time.Local)
		publishTime := t.Unix()
		fmt.Println(fmt.Sprintf("[process]%s创建时间为%s", title, e.Text))

		go rdb.HSet(context.Background(), config.Task.Home, detailUrl, title).Result() //放入缓存

		if time.Now().Unix()-publishTime < config.INTERVAL && !isVisited {
			message := fmt.Sprintf("监测到新的帖子\n标题：%s\n链接：%s\n发布时间：%s", title, detailUrl, e.Text)
			go notification(message) // 触发通知
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Referer", config.Task.Home)
		fmt.Println("[process]VisitDetail:", title)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("[process]error status :", response.StatusCode)
	})

	visited, err := c.HasVisited(detailUrl)
	if err != nil {
		fmt.Println("[process]c.HasVisited err:", err)
		return
	}
	if visited {
		fmt.Println(fmt.Sprintf("[process]%s has visited", detailUrl))
		return
	}
	var count = 0
begin:
	count++
	err = c.Visit(detailUrl)
	if err != nil {
		if count > 3 {
			fmt.Println("[process]c.Visit err:", err)
			return
		}
		time.Sleep(5 * time.Second)
		goto begin
	}
}

func notification(content string) {
	token := wechat.GetAccessToken("ww531df613e9b51972", "iV0r9_rU6PU1-TYQzKTmi5kTEG2RQNvFpQOEcRSvN0g")
	if token == "" {
		fmt.Println("[process]获取token失败")
		return
	}
	wechat.SendMessage(token, content)
}

func initCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.47"),
	)

	c = ip.SetProxy(c)

	return c
}

package process

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"spider_douban/cache"
	"spider_douban/wechat"
	"time"
)

var IsFirst = true

func VisitDetail(url, detailUrl, title string, start int64) {
	c := colly.NewCollector()
	rdb := cache.GetRedisClient()
	defer rdb.Close()

	c.OnHTML(".create-time.color-green", func(e *colly.HTMLElement) {
		fmt.Println(fmt.Sprintf("[%s]创建时间为%s", title, e.Text))
		t, _ := time.ParseInLocation("2006-01-02 03:04:05", e.Text, time.Local)
		publishTime := t.Unix()
		if publishTime > start {
			exists := rdb.HExists(context.Background(), url, detailUrl).Val()

			//新帖子
			if !exists {
				fmt.Println(title)
				go rdb.HSet(context.Background(), url, detailUrl, title).Result() //放入缓存

				if !IsFirst {
					message := fmt.Sprintf("监测到新的帖子\n标题：%s\n链接：%s", title, detailUrl)
					go notification(message) // 触发通知
				}

			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("VisitDetail:", title, "url", detailUrl)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("error status :", response.StatusCode)
	})

	visited, err := c.HasVisited(detailUrl)
	if err != nil {
		fmt.Println("[process]c.HasVisited err:", err)
		return
	}
	if visited {
		fmt.Println(fmt.Sprintf("%s has visited", detailUrl))
		return
	}
	err = c.Visit(detailUrl)
	if err != nil {
		fmt.Println("[process]c.Visit err:", err)
		return
	}
}

func notification(content string) {
	token := wechat.GetAccessToken("ww531df613e9b51972", "iV0r9_rU6PU1-TYQzKTmi5kTEG2RQNvFpQOEcRSvN0g")
	if token == "" {
		fmt.Println("获取token失败")
		return
	}
	wechat.SendMessage(token, content)
}

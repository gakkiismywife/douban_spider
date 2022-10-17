package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly/v2"
	"os"
	"spider_douban/wechat"
	"strings"
	"time"
)

var rdb *redis.Client

var ctx context.Context

var visited bool

var isFirst = true

var group string

func main() {
	flag.StringVar(&group, "group", "", "小组链接")
	flag.Parse()

	if group == "" {
		fmt.Println("请输入需要监控的小组url")
		return
	}

	//初始化redis
	initRedis()

	//删除上次允许的缓存
	rdb.Del(ctx, group)

	//定时
	ticker := time.NewTicker(time.Second * 300)

	c := initCollector(group)

	c.Visit(group)
	for {
		select {
		case <-ticker.C:
			isFirst = false
			visited, _ = c.HasVisited(group)
			if visited {
				c.Init()
			}
			err := c.Visit(group)
			if err != nil {
				fmt.Println("c.Visit err:", err)
				os.Exit(1)
			}
		}
	}
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "120.78.67.238:6379",
		Password: "woaini1!",
		DB:       0,
	})

	ctx = context.Background()
}

func initCollector(url string) *colly.Collector {
	c := colly.NewCollector()

	c.OnHTML("tr td:nth-of-type(1) a", func(e *colly.HTMLElement) {
		//帖子标题
		title := e.Text

		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)

		filter := strings.Contains(title, "【作业】")
		if !filter {
			return
		}

		//链接
		postUrl := e.Attr("href")

		exists := rdb.HExists(ctx, url, postUrl).Val()

		//新帖子
		if !exists {
			fmt.Println(title)
			go rdb.HSet(ctx, url, postUrl, title).Result() //放入缓存

			if !isFirst {
				message := fmt.Sprintf("监测到新的帖子\n标题：%s\n链接：%s", title, postUrl)
				go notification(message) // 触发通知
			}

		}

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Time", time.Now().Format("2006-01-02 03:04:05"), "Visiting", r.URL)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("status:", response.StatusCode)
		fmt.Println("body:", string(response.Body))
		fmt.Println("error:", err)
	})

	return c
}

func notification(content string) {
	token := wechat.GetAccessToken("ww531df613e9b51972", "iV0r9_rU6PU1-TYQzKTmi5kTEG2RQNvFpQOEcRSvN0g")
	if token == "" {
		fmt.Println("获取token失败")
		return
	}
	wechat.SendMessage(token, content)
}

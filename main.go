package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly/v2"
	"math/rand"
	"os"
	"spider_douban/cache"
	"spider_douban/process"
	"strings"
	"time"
)

var rdb *redis.Client

var ctx context.Context

var visited bool

var group string

var start = time.Now().Unix()

func init() {
	flag.StringVar(&group, "group", "", "小组链接")
	flag.Parse()

	if group == "" {
		fmt.Println("请输入需要监控的小组url")
		os.Exit(1)
	}

	rdb = cache.GetRedisClient()
	ctx = context.Background()

	//删除上次允许的缓存
	rdb.Del(ctx, group)
}

func main() {
	//定时
	ticker := time.NewTicker(time.Second * 600)

	c := initCollector()

	err := c.Visit(group)
	if err != nil {
		fmt.Println("[main]c.Visit err:", err)
		return
	}
	for {
		select {
		case <-ticker.C:
			process.IsFirst = false
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

func initCollector() *colly.Collector {
	c := colly.NewCollector()

	c.OnHTML("tr td:nth-of-type(1) a", func(e *colly.HTMLElement) {
		//帖子标题
		title := e.Text
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)

		//过滤标题带作业的
		filter := strings.Contains(title, "【作业】")
		if !filter {
			return
		}

		//链接
		postUrl := e.Attr("href")

		num := time.Duration(rand.Intn(10) + 10)
		time.Sleep(time.Second * num)

		go process.VisitDetail(group, postUrl, title, start)
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

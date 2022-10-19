package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"math/rand"
	"os"
	"spider_douban/config"
	"spider_douban/ip"
	"spider_douban/process"
	"strings"
	"time"
)

//监测的小组url
var group string

func init() {
	flag.StringVar(&group, "group", "", "小组链接")
	flag.Parse()

	if group == "" {
		fmt.Println("请输入需要监控的小组url")
		os.Exit(1)
	}
}

func main() {
	ticker := time.NewTicker(config.INTERVAL * time.Second)

	c := initCollector()
	err := c.Visit(group)
	if err != nil {
		fmt.Println("[main]c.Visit err:", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			var count int8
		again:
			count++
			c = initCollector()
			err := c.Visit(group)
			if err != nil {
				fmt.Println("c.Visit err:", err)
				if count > 3 {
					log.Println("[main] c.Visit limited")
					os.Exit(1)
				}
				time.Sleep(20 * time.Second)
				goto again
			}
		}
	}
}

func initCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.47"),
	)

	c = ip.SetProxy(c)

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

		//随机sleep 3到5秒
		num := time.Duration(rand.Intn(3) + 3)
		time.Sleep(time.Second * num)

		//浏览详情
		go process.VisitDetail(group, postUrl, title)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Time", time.Now().Format("2006-01-02 15:04:05"), "Visiting", r.URL)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("status:", response.StatusCode)
		fmt.Println("body:", string(response.Body))
		fmt.Println("error:", err)
	})

	c.OnResponse(func(response *colly.Response) {
		err := response.Save("response.txt")
		if err != nil {
			return
		}
	})

	return c
}

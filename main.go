package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"math/rand"
	"spider_douban/config"
	"spider_douban/ip"
	"spider_douban/process"
	"strings"
	"time"
)

var c *colly.Collector

var success bool

func main() {
	ticker := time.NewTicker(config.INTERVAL * time.Second)

	run()

	for {
		select {
		case <-ticker.C:
			run()
		}
	}
}

func run() {
	for _, url := range config.Task.Urls {
	again:
		success = true
		c = initCollector()
		err := c.Visit(url)
		c.Wait()
		if err != nil || success == false {
			fmt.Println("[main]c.Visit err:", err)
			i := rand.Intn(20) + 30
			time.Sleep(time.Duration(i) * time.Second)
			goto again
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

		//打印爬取到的帖子标题和时间
		now := time.Now().Format("2006-01-02 15:04:05")
		log.Println(fmt.Sprintf("[main]time:%s,title:%s", now, title))

		//链接
		postUrl := e.Attr("href")

		//随机sleep 3到5秒
		num := time.Duration(rand.Intn(3) + 3)
		time.Sleep(time.Second * num)

		//浏览详情
		go process.VisitDetail(postUrl, title)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[main]Time", time.Now().Format("2006-01-02 15:04:05"), "Visiting", r.URL)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("[main]status:", response.StatusCode)
		fmt.Println("[main]body:", string(response.Body))
		fmt.Println("[main]error:", err)
	})

	c.OnResponse(func(response *colly.Response) {
		body := string(response.Body)
		//判断响应是否正常
		if !strings.Contains(body, "td") {
			log.Println("[main]response body err:", body)
			success = false
		}
	})

	return c
}

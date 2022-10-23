package spider

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"math/rand"
	"spider_douban/cache"
	"spider_douban/config"
	"strings"
	"time"
)

type ListTask struct {
	*Task
	Filters []string
}

func NewListTask(flag, url string, filters []string) *ListTask {
	t := NewTask(flag, url)
	l := &ListTask{
		Task:    nil,
		Filters: filters,
	}
	t.SetRequestCallback(l.requestHandle)
	t.SetHtmlCallback("tr td:nth-of-type(1) a", l.htmlHandle)
	t.SetResponseCallback(l.responseHandle)
	l.Task = t
	return l
}

func (l *ListTask) htmlHandle(e *colly.HTMLElement) {
	//帖子标题
	title := e.Text
	title = strings.Replace(title, " ", "", -1)
	title = strings.Replace(title, "\n", "", -1)
	if len(l.Filters) > 0 {
		for _, filter := range l.Filters {
			if !strings.Contains(title, filter) {
				return
			}
		}
	}

	//打印爬取到的帖子标题和时间
	log.Println(fmt.Sprintf("[%s]%s", l.Flag, title))

	//链接
	postUrl := e.Attr("href")

	rdb := cache.GetRedisClient()
	defer rdb.Close()

	//判断是否请求过
	isVisited := rdb.HExists(context.Background(), config.Task.Home, postUrl).Val()
	if isVisited {
		log.Println(fmt.Sprintf("[%s]%s has visited", l.Flag, title))
		return
	}

	//随机sleep 3到5秒
	num := time.Duration(rand.Intn(3) + 3)
	time.Sleep(time.Second * num)

	//浏览详情
	go func() {
		p := NewPageTask("detail", postUrl, title)
		p.Run(3)
	}()
}

func (l *ListTask) requestHandle(request *colly.Request) {
	log.Println(fmt.Sprintf("[%s]Visiting:%s", l.Flag, l.Url))
}

func (l *ListTask) responseHandle(response *colly.Response) {
	body := string(response.Body)
	//判断响应是否正常
	if !strings.Contains(body, "td") {
		log.Println(fmt.Sprintf("[%s]response body err:%s", l.Flag, body))
		l.State = false
	} else {
		l.State = true
	}
}
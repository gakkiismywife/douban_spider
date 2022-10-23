package spider

import (
	"github.com/gocolly/colly/v2"
	"math/rand"
	"spider_douban/ip"
	"time"
)

type Task struct {
	Url        string
	Header     map[string]string
	Flag       string
	c          *colly.Collector
	OnRequest  colly.RequestCallback
	OnResponse colly.ResponseCallback
	OnHtml     colly.HTMLCallback
	State      bool
}

func NewTask(flag, url string) *Task {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.47"),
	)

	c = ip.SetProxy(c)
	return &Task{
		Url:   url,
		Flag:  flag,
		c:     colly.NewCollector(),
		State: true,
	}
}

// SetHtmlCallback 设置html回调
func (t *Task) SetHtmlCallback(selector string, f colly.HTMLCallback) {
	t.c.OnHTML(selector, f)
}

// SetRequestCallback 设置请求回调
func (t *Task) SetRequestCallback(f colly.RequestCallback) {
	t.c.OnRequest(f)
}

// SetResponseCallback 设置响应回调
func (t *Task) SetResponseCallback(f colly.ResponseCallback) {
	t.c.OnResponse(f)
}

// Run 运行
func (t *Task) Run(times int8) {
	if times > 0 {
		for times > 0 {
			err := t.c.Visit(t.Url)
			t.c.Wait()
			if err != nil || t.State == false {
				time.Sleep(getSleepSecond(10, 20))
				times--
			}
		}
	} else {
	begin:
		err := t.c.Visit(t.Url)
		t.c.Wait()
		if err != nil || t.State == false {
			time.Sleep(getSleepSecond(10, 20))
			goto begin
		}
	}
}

// 生成随机等待时间
func getSleepSecond(min, max int) time.Duration {
	number := rand.Intn(max-min) + min
	return time.Second * time.Duration(number)
}

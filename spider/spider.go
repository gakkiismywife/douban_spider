package spider

import "github.com/gocolly/colly/v2"

type Task struct {
	Url        string
	Header     map[string]string
	c          *colly.Collector
	OnRequest  colly.RequestCallback
	OnError    colly.ErrorCallback
	OnResponse colly.ResponseCallback
}

func NewTask(url string, headers map[string]string) *Task {
	return &Task{
		Url:    url,
		Header: headers,
		c:      colly.NewCollector(),
	}
}

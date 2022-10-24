package main

import (
	"flag"
	"spider_douban/spider"
)

var title string
var url string
var publishTime string

func main() {
	flag.StringVar(&title, "title", "", "标题")
	flag.StringVar(&url, "url", "", "链接")
	flag.StringVar(&publishTime, "time", "", "创建时间")

	flag.Parse()
	if title == "" || url == "" || publishTime == "" {
		return
	}

	spider.Send(title, url, publishTime)
}

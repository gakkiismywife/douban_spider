package main

import (
	"spider_douban/config"
	"spider_douban/spider"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Duration(config.Task.Interval) * time.Second)

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
		l := spider.NewListTask("list", url, []string{"【作业】"})
		l.Run(15)
	}
}

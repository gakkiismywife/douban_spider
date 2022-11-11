package main

import (
    "fmt"
    "os"
    "spider_douban/config"
	"spider_douban/spider"
	"time"
)

func main() {
    fmt.Println(config.Task.FilterWords)
    os.Exit(1)
	go spider.RunSchedule()
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
		l := spider.NewListTask("list", url, config.Task.FilterWords)
		l.Run(15)
	}
}

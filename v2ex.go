package main

import (
	"spider_douban/v2ex"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Duration(3600) * time.Second)

	v2ex.Run()

	for {
		select {
		case <-ticker.C:
			v2ex.Run()
		}
	}
}

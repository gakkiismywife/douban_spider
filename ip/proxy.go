package ip

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"log"
	"net/http"
	"net/url"
)

type ResponseData struct {
	Count int      `json:"count"`
	List  []string `json:"proxy_list"`
}

type ProxyResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data ResponseData `json:"data"`
}

func GetOneIp() string {
	auth := Authorization{SecretId: "945398226373963", SecretKey: "suz1i62o2invo314ouuvdfe84dyf67qc"}
	client := Client{Auth: auth}
	params := map[string]interface{}{"format": "json"}
	ips, err := client.GetDps(1, HmacSha1, params)
	if err != nil {
		log.Println(err)
		return ""
	}
	return ips[0]
}

func SetProxy(c *colly.Collector) *colly.Collector {
	switcher, err := proxy.RoundRobinProxySwitcher("http://t11666754707153:qzgfe0ed@tps163.kdlapi.com:15818")
	if err != nil {
		fmt.Println("[main]proxy.RoundRobinProxySwitcher err", err)
		return c
	}

	u, _ := url.Parse("http://t11666754707153:qzgfe0ed@tps163.kdlapi.com:15818")
	//u.User = url.UserPassword("admin", "o4hwvwob")
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}
	c.WithTransport(transport)
	c.SetProxyFunc(switcher)
	return c
}

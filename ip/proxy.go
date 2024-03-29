package ip

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"log"
	"net/http"
	"net/url"
	"spider_douban/config"
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
	auth := Authorization{SecretId: config.Proxy.Id, SecretKey: config.Proxy.Secret}
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
	if config.Proxy.Tunnel == "" {
		return c
	}
	switcher, err := proxy.RoundRobinProxySwitcher(config.Proxy.Tunnel)
	if err != nil {
		fmt.Println("[main]proxy.RoundRobinProxySwitcher err", err)
		return c
	}

	u, _ := url.Parse(config.Proxy.Tunnel)
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}
	c.WithTransport(transport)
	c.SetProxyFunc(switcher)
	return c
}

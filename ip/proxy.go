package ip

import (
	"github.com/gocolly/colly/v2"
	"log"
	"net"
	"net/http"
	"net/url"
	"spider_douban/config"
	"time"
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
	u, _ := url.Parse(config.Proxy.Tunnel)
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}
	c.WithTransport(transport)
	return c
}

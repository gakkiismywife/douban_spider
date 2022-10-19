package ip

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"io/ioutil"
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

func GetIp() string {
	serverUrl := "https://dps.kdlapi.com/api/getdps/?secret_id=om9zx4xaja6q962r2r52&num=1&signature=oavi2aj83p47qnp1m5hj3lv658&pt=1&format=json&sep=1"
	response, err := http.Get(serverUrl)
	if err != nil {
		fmt.Println("[ip.proxy] http.Get err :", err)
		return ""
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("[ip.proxy] ioutil.ReadAll err :", err)
		return ""
	}

	res := new(ProxyResponse)
	err = json.Unmarshal(all, res)
	if err != nil {
		fmt.Println("[ip.proxy] json.Unmarshal err :", err)
		return ""
	}
	if res.Code != 0 {
		fmt.Println("[ip.proxy] msg :", res.Msg)
		return ""
	}

	if len(res.Data.List) == 0 {
		return ""
	}
	return res.Data.List[0]
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

package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	serverUrl := "https://dps.kdlapi.com/api/getdps/?secret_id=om9zx4xaja6q962r2r52&num=1&signature=o04yhwf64f93qnzh5g9be55sm2&pt=1&format=json&sep=1"
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

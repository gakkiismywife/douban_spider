package ip

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Authorization 用于保存用户SecretId、secretkey，以及计算签名
type Authorization struct {
	SecretId  string
	SecretKey string
}

func TypeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return ""
	}
}

// GetStringToSign 生成签名原文字符串
func (auth Authorization) GetStringToSign(method string, endpoint string, params map[string]interface{}) string {
	s := method + strings.Split(endpoint, ".com")[1] + "?"
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var arr []string
	for _, k := range keys {
		arr = append(arr, k+"="+TypeSwitcher(params[k]))
	}
	queryStr := strings.Join(arr, "&")
	return s + queryStr
}

// SignStr 计算签名串
func (auth Authorization) SignStr(rawStr string) string {
	key := []byte(auth.SecretKey)
	hash := hmac.New(sha1.New, key)
	hash.Write([]byte(rawStr))
	sig := base64.StdEncoding.EncodeToString([]byte(string(hash.Sum(nil))))
	return sig
}

package wechat

type Message interface {
	GenerateMessage(title, url, time string) string
}

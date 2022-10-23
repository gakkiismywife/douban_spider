package spider

import "github.com/gocolly/colly/v2"

type TaskInterface interface {
	htmlHandle(e *colly.HTMLElement)
	requestHandle(request *colly.Request)
	responseHandle(response *colly.Response)
}

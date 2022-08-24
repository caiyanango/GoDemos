package main

import (
	"crawl/regexptest"
	"fmt"
)

func main() {
	/*c := colly.NewCollector()
	c.OnHTML("a", func(element *colly.HTMLElement) {
		element.Request.Visit(element.Attr("href"))
	})
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	c.Visit("http://49.5.6.85:7777")*/
	var start, end int
	fmt.Print("请输入爬虫起始页:")
	fmt.Scanf("%d\r\n", &start)
	fmt.Print("请输入爬虫结束页:")
	fmt.Scanf("%d\r\n", &end)
	regexptest.CrawlPage(start, end)
}

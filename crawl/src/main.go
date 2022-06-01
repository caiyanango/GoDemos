package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	c.OnHTML("a", func(element *colly.HTMLElement) {
		element.Request.Visit(element.Attr("href"))
	})
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	c.Visit("http://49.5.6.85:7777")
}

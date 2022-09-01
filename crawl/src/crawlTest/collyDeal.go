package crawlTest

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
)

func getMovieInfoByColly(info *[]movieInfo, url string) error {
	var moviename, details, strlink []string
	var flag bool
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
	c.DetectCharset = true
	c.OnHTML(".ulink", func(element *colly.HTMLElement) {
		if !strings.Contains(element.Text, "《") {
			return
		}
		name := strings.Split(strings.Split(element.Text, "《")[1], "》")
		if name == nil {
			name = strings.Split(strings.Split(element.Text, "《")[1], "】")
		}
		moviename = append(moviename, name[0])
		flag = false
		strlink = append(strlink, "/")
		if link := element.Attr("href"); link != "" {
			c.Visit(serveraddr + link)
		}
	})
	c.OnHTML("a[href^=magnet],a[href^=ftp]", func(element *colly.HTMLElement) {
		if !flag {
			link := element.Attr("href")
			if !strings.HasPrefix(link, "magnet") && !strings.HasPrefix(link, "ftp") {
				link = "/"
			}
			strlink = append(strlink[:len(strlink)-1], link)
			flag = true
		}
	})
	c.OnHTML("td[colspan=\"2\"][style=\"padding-left:3px\"]", func(element *colly.HTMLElement) {
		details = append(details, element.Text)
	})
	c.Visit(url)
	if len(moviename) != len(details) || len(strlink) != len(moviename) {
		return fmt.Errorf("movie count [%d] details count [%d] strlink count [%d] is not equal on %s page", len(moviename), len(details), len(strlink), url)
	}
	for i := 0; i < len(moviename); i++ {
		var movieinfoElement movieInfo
		fillMovieElementInfo(&movieinfoElement, moviename[i], details[i], strlink[i])
		*info = append(*info, movieinfoElement)
	}
	return nil
}

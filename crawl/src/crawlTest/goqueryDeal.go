package crawlTest

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func getUrlDom(url string) (*goquery.Document, error) {
	result, err := getUrlContent(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(result))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dom, nil
}

func getMovieInfoByGoquery(info *[]movieInfo, url string) error {
	dom, err := getUrlDom(url)
	if err != nil {
		log.Println(err)
		return err
	}
	var moviename, details, strlink []string
	dom.Find(".ulink").Each(func(i int, selection *goquery.Selection) {
		if !strings.Contains(selection.Text(), "《") {
			return
		}
		name := strings.Split(strings.Split(selection.Text(), "《")[1], "》")
		if name == nil {
			name = strings.Split(strings.Split(selection.Text(), "《")[1], "】")
		}
		moviename = append(moviename, name[0])
		link, _ := selection.Attr("href")
		dom, err := getUrlDom(serveraddr + link)
		if err != nil {
			log.Println(err)
			link = "/"
		} else {
			dom.Find("a[href^=magnet],a[href^=ftp]").Each(func(i int, selection *goquery.Selection) {
				link, _ = selection.Attr("href")
			})
			if !strings.HasPrefix(link, "magnet") && !strings.HasPrefix(link, "ftp") {
				link = "/"
			}
		}
		strlink = append(strlink, link)
	})
	dom.Find("td[colspan=\"2\"][style=\"padding-left:3px\"]").Each(func(i int, selection *goquery.Selection) {
		details = append(details, selection.Text())
	})
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

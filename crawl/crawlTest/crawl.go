package crawlTest

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

const serveraddr = "https://m.dytt8.net"

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

type movieInfo struct {
	moviename   string
	year        string
	from        string
	kind        string
	language    string
	subtitle    string
	releaseDate string
	score       string
	link        string
}

func CrawlPage(start, end int, mode *string) {
	file, sheet, headers, err := createFile()
	if err != nil {
		log.Fatal(err)
	}
	countOfWorkGoroutine := end - start + 1
	wg.Add(countOfWorkGoroutine)
	for i := start; i <= end; i++ {
		go func(i int) {
			var err error
			defer func() {
				if err != nil {
					log.Printf("第%d页爬取失败, %v\n", i, err)
				} else {
					fmt.Printf("第%d页已爬取完成\n", i)
				}
				wg.Done()
			}()
			movieinfo := []movieInfo{}
			url := serveraddr + "/html/gndy/dyzz/list_23_" + strconv.Itoa(i) + ".html"
			switch *mode {
			case "regext":
				err = getMovieInfo(&movieinfo, url)
				if err != nil {
					return
				}
			case "goquery":
				err = getMovieInfoByGoquery(&movieinfo, url)
				if err != nil {
					return
				}
			case "colly":
				err = getMovieInfoByColly(&movieinfo, url)
				if err != nil {
					return
				}
			default:
				log.Panicf("不支持的模式: %s\n", *mode)
			}
			data := make(map[string]string)
			mu.Lock()
			for _, info := range movieinfo {
				data["Name"] = info.moviename
				data["Year"] = info.year
				data["From"] = info.from
				data["Kind"] = info.kind
				data["Language"] = info.language
				data["Subtitle"] = info.subtitle
				data["ReleaseDate"] = info.releaseDate
				data["Score"] = info.score
				data["Link"] = info.link
				row := sheet.AddRow()
				row.SetHeightCM(0.8)
				for _, field := range *headers {
					row.AddCell().Value = data[field.Field]
				}
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	err = saveFile(file, "电影天堂.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("export success")
}

package regexptest

import (
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

type HeaderColumn struct {
	Field string // 字段，数据映射到的数据字段名
	Title string // 标题，表格中的列名称
}

func convertString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
func getUrlContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := fmt.Errorf("visit %s error, http code is %d", url, resp.StatusCode)
		return "", err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	str := string(result)
	if strings.Contains(str, "charset=gb2312") {
		str = convertString(str, "gbk", "utf-8")
	}
	return str, nil
}

func createFile() (*xlsx.File, *xlsx.Sheet, *[]*HeaderColumn, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("MOVIESINFO") //表实例
	if err != nil {
		return nil, nil, nil, err
	}
	headers := []*HeaderColumn{
		{Field: "Name", Title: "电影名称"},
		{Field: "Year", Title: "年代"},
		{Field: "From", Title: "产地"},
		{Field: "Kind", Title: "类别"},
		{Field: "Language", Title: "语言"},
		{Field: "Subtitle", Title: "字幕"},
		{Field: "ReleaseDate", Title: "上映日期"},
		{Field: "Score", Title: "IMDb评分"},
		{Field: "Link", Title: "磁力链接"},
	}
	style := map[string]float64{
		"Name":        2.0,
		"Year":        2.0,
		"From":        2.0,
		"Kind":        2.0,
		"Language":    2.0,
		"Subtitle":    2.0,
		"ReleaseDate": 2.0,
		"Score":       2.0,
		"Link":        2.0,
	}
	sheet, _ = SetHeader(sheet, headers, style)
	return file, sheet, &headers, nil
}

func saveFile(file *xlsx.File, name string) error {
	err := file.Save(name)
	if err != nil {
		return err
	}
	return nil
}

func parseStr(src string, sep ...string) string {
	var result string
	for _, s := range sep {
		if strings.Contains(src, s) {
			result = strings.TrimSpace(strings.Split(src, s)[1])
			break
		}
	}
	return result
}

func CrawlPage(start, end int) {
	file, sheet, headers, err := createFile()
	if err != nil {
		log.Fatal(err)
	}
	countOfWorkGoroutine := end - start + 1
	wg.Add(countOfWorkGoroutine)
	for i := start; i <= end; i++ {
		go func(i int) {
			defer func() {
				fmt.Printf("第%d页已爬取完成\n", i)
				wg.Done()
			}()
			movieinfo := []movieInfo{}
			url := serveraddr + "/html/gndy/dyzz/list_23_" + strconv.Itoa(i) + ".html"
			result, err := getUrlContent(url)
			if err != nil {
				log.Println(err)
				return
			}
			reg := regexp.MustCompile(`<a href="(.+?)" class="ulink">.*?[《【](.+?)[》】].*?</a>`)
			movie := reg.FindAllStringSubmatch(result, -1)
			reg = regexp.MustCompile(`<td colspan="2" style="padding-left:3px">(?s)(.*?)</td>`)
			details := reg.FindAllStringSubmatch(result, -1)
			if len(movie) != len(details) {
				fmt.Println("movie", len(movie), "details", len(details), i)
				return
			}
			for infoIdx := 0; infoIdx < len(movie); infoIdx++ {
				var movieinfoElement movieInfo
				movieinfoElement.moviename = movie[infoIdx][2]
				strKind := strings.Split(details[infoIdx][1], "◎")
				for _, v := range strKind {
					if strings.Contains(v, "年 代") || strings.Contains(v, "年　　代") {
						movieinfoElement.year = parseStr(v, "年 代", "年　　代")
					} else if strings.Contains(v, "产 地") || strings.Contains(v, "国 家") || strings.Contains(v, "国　　家") || strings.Contains(v, "地 区") {
						movieinfoElement.from = parseStr(v, "产 地", "国 家", "国　　家", "地 区")
					} else if strings.Contains(v, "类 别") || strings.Contains(v, "类　　别") {
						movieinfoElement.kind = parseStr(v, "类 别", "类　　别")
					} else if strings.Contains(v, "语 言") || strings.Contains(v, "语　　言") {
						movieinfoElement.language = parseStr(v, "语 言", "语　　言")
					} else if strings.Contains(v, "字 幕") || strings.Contains(v, "字　　幕") {
						movieinfoElement.subtitle = parseStr(v, "字 幕", "字　　幕")
					} else if strings.Contains(v, "上映日期") {
						movieinfoElement.releaseDate = parseStr(v, "上映日期")
					} else if strings.Contains(v, "IMDb评分") || strings.Contains(v, "IMDB评分") {
						movieinfoElement.score = parseStr(v, "IMDb评分", "IMDB评分")
					}
				}
				result, err := getUrlContent(serveraddr + movie[infoIdx][1])
				if err != nil {
					log.Println(err)
					movieinfoElement.link = "/"
				} else {
					reg := regexp.MustCompile(`href="(magnet.+?)"|href="(ftp.+?)"|>(ftp.+?)</a>`)
					link := reg.FindAllStringSubmatch(result, 1)
					if len(link) != 0 {
						if link[0][1] != "" && strings.HasPrefix(link[0][1], "magnet") {
							movieinfoElement.link = link[0][1]
						} else if link[0][2] != "" && strings.HasPrefix(link[0][2], "ftp") {
							movieinfoElement.link = link[0][2]
						} else if link[0][3] != "" && strings.HasPrefix(link[0][3], "ftp") {
							movieinfoElement.link = link[0][3]
						}
					} else {
						movieinfoElement.link = "/"
					}
				}
				movieinfo = append(movieinfo, movieinfoElement)
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

func SetHeader(sheet *xlsx.Sheet, header []*HeaderColumn, width map[string]float64) (*xlsx.Sheet, error) {
	if len(header) == 0 {
		return nil, errors.New("Excel.SetHeader 错误: 表头不能为空")
	}

	// 表头样式
	style := xlsx.NewStyle()

	font := xlsx.DefaultFont()
	font.Bold = true

	alignment := xlsx.DefaultAlignment()
	alignment.Vertical = "center"

	style.Font = *font
	style.Alignment = *alignment

	style.ApplyFont = true
	style.ApplyAlignment = true

	// 设置表头字段
	row := sheet.AddRow()
	row.SetHeightCM(1.0)
	row_w := make([]string, 0)
	for _, column := range header {
		row_w = append(row_w, column.Field)
		cell := row.AddCell()
		cell.Value = column.Title
		cell.SetStyle(style) //设置单元样式
	}

	// 表格列，宽度
	if len(row_w) > 0 {
		for k, v := range row_w {
			if width[v] > 0.0 {
				sheet.SetColWidth(k, k, width[v]*10)
			}
		}
	}
	return sheet, nil
}

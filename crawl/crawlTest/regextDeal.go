package crawlTest

import (
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func convertString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func getResponse(url string) (*http.Response, error) {
	client := http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getUrlContent(url string) (string, error) {
	resp, err := getResponse(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := fmt.Errorf("visit %s error, http code is %d", url, resp.StatusCode)
		return "", err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	str := string(result)
	if strings.Contains(str, "charset=gb2312") {
		str = convertString(str, "gbk", "utf-8")
	}
	return str, nil
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

func fillMovieElementInfo(movieinfoElement *movieInfo, moviename, details, link string) {
	movieinfoElement.moviename = moviename
	strKind := strings.Split(details, "◎")
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
	movieinfoElement.link = link
}

func getMovieInfo(info *[]movieInfo, url string) error {
	result, err := getUrlContent(url)
	if err != nil {
		log.Println(err)
		return err
	}
	reg := regexp.MustCompile(`<a href="(.+?)" class="ulink">.*?[《【](.+?)[》】].*?</a>`)
	movie := reg.FindAllStringSubmatch(result, -1)
	reg = regexp.MustCompile(`<td colspan="2" style="padding-left:3px">(?s)(.*?)</td>`)
	details := reg.FindAllStringSubmatch(result, -1)
	if len(movie) != len(details) {
		return fmt.Errorf("movie count [%d] is not equal details count [%d] on %s page", len(movie), len(details), url)
	}
	for infoIdx := 0; infoIdx < len(movie); infoIdx++ {
		var movieinfoElement movieInfo
		var strlink string
		result, err := getUrlContent(serveraddr + movie[infoIdx][1])
		if err != nil {
			log.Println(err)
			strlink = "/"
		} else {
			reg := regexp.MustCompile(`href="(magnet.+?)"|href="(ftp.+?)"|>(ftp.+?)</a>`)
			link := reg.FindAllStringSubmatch(result, 1)
			if len(link) != 0 {
				if link[0][1] != "" && strings.HasPrefix(link[0][1], "magnet") {
					strlink = link[0][1]
				} else if link[0][2] != "" && strings.HasPrefix(link[0][2], "ftp") {
					strlink = link[0][2]
				} else if link[0][3] != "" && strings.HasPrefix(link[0][3], "ftp") {
					strlink = link[0][3]
				}
			} else {
				strlink = "/"
			}
		}
		fillMovieElementInfo(&movieinfoElement, movie[infoIdx][2], details[infoIdx][1], strlink)
		*info = append(*info, movieinfoElement)
	}
	return nil
}

package main

import (
	"bufio"
	"crawl/crawlTest"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	crawlMode := flag.String("mode", "regext", "specify the crawler method")
	flag.Parse()
	fmt.Printf("使用的爬虫模式: %s\n", *crawlMode)
	var start, end int
	fmt.Print("请输入爬虫起始页:")
	fmt.Scanln(&start)
	fmt.Print("请输入爬虫结束页:")
	fmt.Scanln(&end)
	timeStart := time.Now()
	crawlTest.CrawlPage(start, end, crawlMode)
	fmt.Printf("用时%.2f秒\n", time.Since(timeStart).Seconds())
	StopUntil("请按任意键退出", "", false)
	//StopUntil("请输入exit退出", "exit", true)
	fmt.Print("\n")
}

func StopUntil(hint string, trigger string, repeat bool) error {
	pressLen := len([]rune(trigger))
	if trigger == "" {
		pressLen = 1
	}
	fd := int(os.Stdin.Fd())
	if hint != "" {
		fmt.Print(hint)
	}
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)
	inputReader := bufio.NewReader(os.Stdin)
	i := 0
	for {
		b, _, err := inputReader.ReadRune()
		if err != nil {
			return err
		}
		if trigger == "" {
			break
		}
		if b == []rune(trigger)[i] {
			i++
			if i == pressLen {
				break
			}
			continue
		}
		i = 0
		if hint != "" && repeat {
			fmt.Print("\n")
			fmt.Print(hint)
		}
	}
	return nil
}

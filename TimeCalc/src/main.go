package main

import (
	"fmt"
	"time"
)

func main() {
	var year, month, day int
	var count int
	for {
		fmt.Println("请输入候选人的毕业时间(某年某月某日):")
		fmt.Scanf("%d年%d月%d日\r\n", &year, &month, &day)
		loc, _ := time.LoadLocation("Asia/Shanghai")
		preTimePoint := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
		workTime := int(time.Since(preTimePoint).Hours())
		fmt.Printf("候选人的工作年限是%d年%d个月\n", workTime/24/365, workTime/24%365/30)
		for {
			select {
			case <-time.Tick(time.Second):
				count++
			}
			if count == 10 {
				count = 0
				break
			}
		}
	}
}

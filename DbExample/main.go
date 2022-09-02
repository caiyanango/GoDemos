package main

import (
	"DbExample/mongodb_example"
	"DbExample/mysql_example"
	"DbExample/redis_example"
	"fmt"
)

func main() {
	var dbKind int
	for {
		fmt.Println("1、MySQL   2、MongoDB   3、Redis")
		fmt.Print("Select which database to operate:")
		fmt.Scan(&dbKind)
		if dbKind == 1 {
			fmt.Println("MySQL Test.")
			var opKind int
			fmt.Println("1、RawSQL   2、Gorm   3、Xorm")
			fmt.Print("Select which option to operate:")
			fmt.Scan(&opKind)
			if opKind == 1 {
				mysql_example.Connect()
			} else if opKind == 2 {
				mysql_example.ConnectByGorm()
			} else {
				mysql_example.ConnectByXorm()
			}
		} else if dbKind == 2 {
			fmt.Println("MongoDB Test.")
			mongodb_example.Connect()
		} else if dbKind == 3 {
			fmt.Println("Redis Test.")
			redis_example.Connect()
		}
	}
}

package main

import (
	"DbExample/mongodb_example"
	"DbExample/mysql_example"
	"fmt"
)

func main() {
	var dbKind int
	for {
		fmt.Println("1、MySQL   2、MongoDB")
		fmt.Print("Select which database to operate:")
		fmt.Scan(&dbKind)
		if dbKind == 1 {
			fmt.Println("MySQL Test.")
			mysql_example.Connect()
		} else if dbKind == 2 {
			fmt.Println("MongoDB Test.")
			mongodb_example.Connect()
		}
	}
}

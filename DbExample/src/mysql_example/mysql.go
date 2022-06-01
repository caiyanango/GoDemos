package mysql_example

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
)

const (
	insert = iota + 1
	delete
	update
	retrieve
	disconect
)

type userInfo struct {
	id       int
	username string
	password string
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func processData(db *sql.DB, cmd int, data ...any) error {
	var sqlStr string
	switch cmd {
	case insert:
		sqlStr = "INSERT INTO user_tb1 (username,password) VALUES (?,?)"
	case delete:
		sqlStr = "DELETE FROM user_tb1 WHERE id=?"
	case update:
		sqlStr = fmt.Sprintf("UPDATE user_tb1 SET username=?,password=? WHERE id=%s", data[0])
		data = data[1:]
	}
	_, err := db.Exec(sqlStr, data...)
	if err != nil {
		return err
	}
	return nil
}

func retrieveData(db *sql.DB, condition string) ([]userInfo, error) {
	result := []userInfo{}
	var u userInfo
	var err error
	if strings.EqualFold(condition, "all") {
		condition = ""
	}
	sqlStr := fmt.Sprintf("SELECT * FROM user_tb1 %s", condition)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&u.id, &u.username, &u.password)
		if err != nil {
			fmt.Printf("scan failed,%s\n", err)
			continue
		}
		result = append(result, u)
	}
	return result, nil
}

func Connect() {
	db, err := sql.Open("mysql", "root:woaichx199298@tcp(192.168.9.53:3306)/go_db")
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	fmt.Println("Connect database success.")
	defer func() {
		db.Close()
		fmt.Println("Disconnect database success.")
	}()
	var cmd int
	for {
		fmt.Println("Please select your operation:")
		fmt.Println("1、Insert  2、Delete  3、Update  4、Retrieve  5、Disconnect")
		fmt.Scan(&cmd)
		switch cmd {
		case insert:
			{
				var username, password string
				fmt.Print("Enter username:")
				fmt.Scan(&username)
				fmt.Print("Enter password:")
				fmt.Scan(&password)
				err := processData(db, cmd, username, password)
				if err != nil {
					fmt.Printf("insert data failed, %s\n", err)
				} else {
					fmt.Println("insert data success.")
				}
			}
		case retrieve:
			{
				var condition string
				fmt.Print("Enter condition:")
				reader := bufio.NewReader(os.Stdin)
				for {
					condition, _ = reader.ReadString('\n')
					if condition == "\r\n" || condition == "\n" {
						continue
					}
					break
				}
				if strings.Contains(condition, "\r\n") {
					condition = condition[0 : len(condition)-2]
				} else if strings.Contains(condition, "\n") {
					condition = condition[0 : len(condition)-1]
				}
				result, err := retrieveData(db, condition)
				if err != nil {
					fmt.Printf("retrieve data failed, %s\n", err)
					break
				}
				fmt.Printf("%-5s      %-16s      %-16s\n", "id", "username", "password")
				for _, info := range result {
					fmt.Printf("%-5d      %-16s      %-16s\n", info.id, info.username, info.password)
				}
			}
		case update:
			{
				var id, username, password string
				fmt.Print("Enter id you want to update:")
				fmt.Scan(&id)
				fmt.Print("Enter new username:")
				fmt.Scan(&username)
				fmt.Print("Enter new password:")
				fmt.Scan(&password)
				err := processData(db, cmd, id, username, password)
				if err != nil {
					fmt.Printf("update data failed, %s\n", err)
				} else {
					fmt.Println("update data success.")
				}
			}
		case delete:
			{
				var deleteID int
				fmt.Print("Enter id you want to delete:")
				fmt.Scan(&deleteID)
				err := processData(db, cmd, deleteID)
				if err != nil {
					fmt.Printf("delete data failed, %s\n", err)
				} else {
					fmt.Println("delete data success.")
				}
			}
		case disconect:
			return
		default:
			fmt.Println("Unsupport option.")
		}
	}
}

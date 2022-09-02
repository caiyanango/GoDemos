package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MGRDB manage_db

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type manage_db struct {
	db *gorm.DB
}

func init() {
	dsn := "root:123456@tcp(172.16.9.227:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	MGRDB.db = db
}
func (mgr *manage_db) Adduser(user *User) {
	mgr.db.Create(user)
}

func (mgr *manage_db) Getuser(name string, user *User) {
	mgr.db.Where("username = ?", name).First(user)
}

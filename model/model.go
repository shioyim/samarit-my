package model

import (
	"log"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/io"
	"github.com/jinzhu/gorm"
	// for db SQL
	"github.com/shioyim/samarit-my/config"
)

var (
	// DB Database
	DB     *gorm.DB
	dbType = "mysql"
	dbURL  = config.String("dburl")
	username = config.String("administrator")
	password = config.String("password")
)



func init() {
	io.Register((*User)(nil), "User", "json")
	//io.Register((*Other)(nil), "Other", "json")
	var err error  
	DB, err = gorm.Open(dbType, dbURL)

	if err != nil {
		log.Fatalln("Connect to database error:", err)
	}
    
    //DB.AutoMigrate(&User{},&Other{})
	DB.AutoMigrate(&User{},)
	users := []User{}
	DB.Find(&users)
	if len(users) == 0 {
		admin := User{
			Username: username,
			Password: password,
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	DB.LogMode(false)
	go ping()
}

func ping() {
	for {
		if err := DB.Exec("SELECT 1").Error; err != nil {
			log.Println("Database ping error:", err)
			if DB, err = gorm.Open(strings.ToLower(dbType), dbURL); err != nil {
				log.Println("Retry connect to database error:", err)
			}
		}
		time.Sleep(time.Minute)
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	return gorm.Open(strings.ToLower(dbType), dbURL)
}

